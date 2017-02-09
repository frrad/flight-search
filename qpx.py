import requests
import json


class qpx:

    def __init__(self, api_key_path):
        with open(api_key_path, 'r') as f:
            self.__key = f.read().strip()

    def search(self, payload):
        headers = {'content-type': 'application/json'}
        key_suffix = '?key=' + self.__key
        r = requests.post(
            'https://www.googleapis.com/qpxExpress/v1/trips/search' + key_suffix,
            data=json.dumps(payload),
            headers=headers
        )
        return r.json()
