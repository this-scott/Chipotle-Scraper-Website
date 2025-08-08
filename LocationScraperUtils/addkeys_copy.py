import json
import requests
import os

#loading environment variables from .env file
def load_env_var(var_name, default=None):
    try:
        with open(".env", "r") as env_file:
            for line in env_file:
                key, value = line.strip().split("=", 1)
                if key == var_name:
                    return value
    except FileNotFoundError:
        pass
    return default

#API_KEY = load_env_var("google_api")  # Read API key from .env file

def get_lat_long_google(address, city, state, api_key):
    base_url = "https://maps.googleapis.com/maps/api/geocode/json"
    query = f"{address}, {city}, {state}, USA"
    params = {"address": query, "key": api_key}
    response = requests.get(base_url, params=params)
    print(response.headers)
    print(response.content)
    if response.status_code == 200 and response.json().get("results"):
        location = response.json()["results"][0]["geometry"]["location"]
        return float(location["lat"]), float(location["lng"])
    return None, None  # Return None if geocoding fails

#load data from file
with open("output.json", "r") as file:
    data = json.load(file)

#adding 'id' field and 'lat', 'long' fields to each entry
for i, entry in enumerate(data, start=1):
    entry["id"] = i  # Assign a unique numerical ID
    lat, long = get_lat_long_google(entry["address"], entry["city"], entry["state"], API_KEY)
    entry["latitude"] = lat
    entry["longitude"] = long

#saving modified data back to a file
with open("updated_data.json", "w") as file:
    json.dump(data, file, indent=2)
