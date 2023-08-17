import psutil
import platform
import json
import requests
import time
import subprocess


def get_updates_windows():
    updates = {
        "Installed": [],
        "Pending": [],
        "Failed": []
    }
    command = 'powershell "Get-HotFix | Format-List -Property Description,HotfixID,InstalledOn"'
    result = subprocess.run(
        command, capture_output=True, text=True, shell=True)
    updates["Installed"] = result.stdout.splitlines()
    return updates


def get_updates_mac():
    updates = {
        "Pending": []
    }
    command = "softwareupdate -l"
    result = subprocess.run(
        command, capture_output=True, text=True, shell=True)
    updates["Pending"] = result.stdout.splitlines()
    return updates


def get_updates_linux():
    updates = {
        "Pending": []
    }
    command = "apt list --upgradeable"
    result = subprocess.run(command.split(), capture_output=True, text=True)
    updates["Pending"] = result.stdout.splitlines()
    command = "grep 'install ' /var/log/dpkg.log"
    result = subprocess.run(command.split(), capture_output=True, text=True)
    updates["Installed"] = result.stdout.splitlines()
    return updates


def get_nvidia_gpu():
    try:
        result = subprocess.check_output(
            ["nvidia-smi", "--query-gpu=gpu_name", "--format=csv,noheader,nounits"]).decode()
        return result.strip()
    except:
        return "NVIDIA GPU information not available or not an NVIDIA GPU."

#  Motherboard info for windows


def get_motherboard_windows():
    try:
        result = subprocess.check_output(
            ["wmic", "baseboard", "get", "product"]).decode()
        lines = result.strip().split('\n')
        if len(lines) > 1:
            return lines[1]
    except Exception as e:
        pass
    return "Motherboard information not available."

#  Motherboard info for linux


def get_motherboard_linux():
    try:
        result = subprocess.check_output(
            ["sudo", "dmidecode", "-t", "2"]).decode()
        for line in result.split('\n'):
            if 'Product Name:' in line:
                return line.split(":")[1].strip()
    except Exception as e:
        pass
    return "Motherboard information not available or requires elevated permissions."

#  Motherboard info for mac OS


def get_motherboard_mac():
    try:
        result = subprocess.check_output(
            ["system_profiler", "SPHardwareDataType"]).decode()
        for line in result.split('\n'):
            if 'Model Name' in line:
                return line.split(":")[1].strip()
    except Exception as e:
        pass
    return "Motherboard information not available."


def gather_system_info():
    data = {
        'OS': None,
        'RAM': {},
        'CPU': {},
        'Disks': [],
        'Battery': {},
        'Boot_Time': None,
        'System_Uptime': None,
        'Network': {},
        'Updates': {},
        'GPU': get_nvidia_gpu()

    }

    # OS Info
    os_info = platform.system()
    if os_info == "Windows":
        data["Motherboard"] = get_motherboard_windows()
    elif os_info == "Linux":
        data["Motherboard"] = get_motherboard_linux()
    elif os_info == "Darwin":  # macOS
        data["Motherboard"] = get_motherboard_mac()
    else:
        data["Motherboard"] = "Unknown OS, motherboard info not available."

    os_version = platform.version()
    os_release = platform.release()
    data["OS"] = f"{os_info} - Version: {os_version} - Release: {os_release}"

    # RAM Info
    virtual_memory = psutil.virtual_memory()
    data["RAM"] = {
        "Total": f"{virtual_memory.total / (1024 ** 3):.2f} GB",
        "Used": f"{virtual_memory.used / (1024 ** 3):.2f} GB",
        "Free": f"{virtual_memory.available / (1024 ** 3):.2f} GB",
        "Percentage": f"{virtual_memory.percent}%"
    }

    # CPU Info
    try:
        cpu_freq = f"{psutil.cpu_freq().current} MHz"
    except Exception:
        cpu_freq = "Unavailable"

    data["CPU"] = {
        "Cores": psutil.cpu_count(logical=False),
        "Percentage": psutil.cpu_percent(interval=1),
        "Frequency": cpu_freq
    }

    # Disk Info
    partitions = psutil.disk_partitions()
    for partition in partitions:
        usage = psutil.disk_usage(partition.mountpoint)
        data["Disks"].append({
            "Device": partition.device,
            "Total": f"{usage.total / (1024 ** 3):.2f} GB",
            "Used": f"{usage.used / (1024 ** 3):.2f} GB",
            "Free": f"{usage.free / (1024 ** 3):.2f} GB"
        })

    # Battery Info
    if hasattr(psutil, "sensors_battery"):
        battery = psutil.sensors_battery()
        if battery:
            data["Battery"] = {
                "Percentage": f"{battery.percent}%",
                "Time Left": str(battery.secsleft),
                "Plugged In": battery.power_plugged
            }

    # Boot Time and System Uptime
    boot_time = psutil.boot_time()
    current_time = time.time()
    data["Boot_Time"] = boot_time
    data["System_Uptime"] = current_time - boot_time

    # Network Info
    data["Network"] = {
        "Total_Bytes_Sent": psutil.net_io_counters().bytes_sent,
        "Total_Bytes_Received": psutil.net_io_counters().bytes_recv
    }

    # Fetching update information based on the OS
    if os_info == "Windows":
        data["Updates"] = get_updates_windows()
    elif os_info == "Darwin":  # Darwin denotes macOS
        data["Updates"] = get_updates_mac()
    elif os_info == "Linux":
        data["Updates"] = get_updates_linux()
    else:
        data["Updates"] = {"message": "Updates info not available for this OS"}

    return data


def write_data_to_file(data, filename="system_info.txt"):
    with open(filename, 'w') as file:
        file.write(json.dumps(data, indent=4))


def send_data_to_server(data):
    # Endpoint to receive the data
    url = 'YOUR_SERVER_ENDPOINT'
    headers = {'Content-type': 'application/json'}

    # Convert the dictionary to a JSON string
    json_data = json.dumps(data)

    response = requests.post(url, data=json_data, headers=headers)
    return response.text


if __name__ == '__main__':
    system_data = gather_system_info()
    write_data_to_file(system_data)
    server_response = send_data_to_server(system_data)
    print(server_response)  # To see server response
