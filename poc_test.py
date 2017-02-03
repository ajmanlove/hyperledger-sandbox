import os
import os.path
from subprocess import check_call
import sys
from hyperledger.client import Client

import requests
import json


import configparser

setup_hl_creds = ("test_user0", "MS9qrN8hFjlE")
insurer1_hl_creds = ("insurer1", "iF6bJY5sWumg")
reinsurer1_hl_creds = ("reinsurer1", "8jUZV34kUBcN")
reinsurer2_hl_creds = ("reinsurer2", "RmMio3LCk37J")
insurer2_hl_creds = ("insurer2", "1wtM0CXjIVXr")
reinsurer3_hl_creds = ("reinsurer3", "atjQRL2S6FJx")

def system_test():
    if not os.path.exists(".setup.ini"):
        print("poc_test requires the .setup.ini file written by setup_poc.py")
        exit(1)

    config = configparser.ConfigParser()
    config.read('.setup.ini')
    print(config.sections())

    c = Client(base_url="http://127.0.0.1:7050")

    submit(config, c)

    print("Not implemented")
    exit(1)

def submit(config, client):
    data = {
        "jsonrpc": "2.0",
        "method": "invoke",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config['CHAINCODE']['reinsurance_request']
            },
            "ctorMsg": {
                "function": "submit",
                "args": [
                    "reinsurer1,reinsurer2", "2e1b1b0cb7bfce4cf47706752a234f29", "http://mybucket.s3-website-us-east-1.amazonaws.com/", "some excel contract text here", "CREATE ABSTRACT TABLE insuredItem (foo INT);", "1"
                ]
            },
            "secureContext": "insurer1",
            "attributes": ["enrollmentId"]
        },
        "id": 1
    }

    post(data)


    print("submit()")

def propose():
    print("propose()")

def counter():
    print("counter()")

def accept():
    print("accept()")

def reject():
    print("reject()")

def post(data):
    data_json = json.dumps(data)
    headers = {'Content-type': 'application/json'}
    response = requests.post("http://localhost:7050/chaincode", data=data_json, headers=headers)

    print("RESPONSE : ", response)

    if response.status_code != 200:
        print("Unexpected status code in registrar " + response.status_code)
        exit(1)

    print("Response JSON " + response.text)

def init_hyperledger():
    print("Initializing hyperledger environment...")
    register_hl_user(setup_hl_creds[0], setup_hl_creds[1])

    c = Client(base_url="http://127.0.0.1:7050")
    # enroll_cc_name = deploy_chaincode(
    #     c, setup_hl_creds[0], "https://github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/enrollment_service", []
    # )

    asset_cc_name = deploy_chaincode(
        c, setup_hl_creds[0], "https://github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/asset_management", []
    )

    print("Asset chaincode name is " + asset_cc_name)

    request_cc_name = deploy_chaincode(
        c, setup_hl_creds[0],
        "https://github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/reinsurance_request",
        [asset_cc_name]
    )

    print("Request chaincode name is " + request_cc_name)

    proposal_cc_name = deploy_chaincode(
        c, setup_hl_creds[0],
        "https://github.com/ajmanlove/hyperledger-sandbox/reinsurance_poc/reinsurance_proposal",
        [asset_cc_name]
    )

    print("Enrolling test users...")
    register_hl_user(insurer1_hl_creds[0], insurer1_hl_creds[1])
    register_hl_user(reinsurer1_hl_creds[0], reinsurer1_hl_creds[1])
    register_hl_user(reinsurer2_hl_creds[0], reinsurer2_hl_creds[1])
    register_hl_user(insurer2_hl_creds[0], insurer2_hl_creds[1])
    register_hl_user(reinsurer3_hl_creds[0], reinsurer3_hl_creds[1])

    register_cc(asset_cc_name, setup_hl_creds[0], request_cc_name, "reinsurance_request")
    register_cc(asset_cc_name, setup_hl_creds[0], proposal_cc_name, "reinsurance_proposal")

    print("")
    print("-----------------------------------------------------------")
    print("asset_management chaincode name: ", asset_cc_name)
    print("reinsurance_request chaincode name: ", request_cc_name)
    print("reinsurance_proposal chaincode name: ", proposal_cc_name)
    print("-----------------------------------------------------------")
    print("")
    print("Init of hyperledger environment COMPLETE")

def register_cc(am_name, user, cc_name, identifier):
    print("Registering chaincode {} as {}".format(cc_name, identifier))
    data = {
      "jsonrpc": "2.0",
      "method": "invoke",
      "params": {
        "type": 1,
        "chaincodeID": {
          "name": am_name
        },
        "ctorMsg": {
          "function": "register_chaincode",
          "args": [cc_name, identifier]
        },
        "secureContext": user,
        "attributes": ["enrollmentId", "contact"]
      },
      "id": 3
    }

    data_json = json.dumps(data)
    headers = {'Content-type': 'application/json'}
    response = requests.post("http://localhost:7050/chaincode", data=data_json, headers=headers)

    print("RESPONSE : ", response)

    if response.status_code != 200:
        print("Unexpected status code in registrar " + response.status_code)
        exit(1)

    print("Response JSON " + response.text)


def deploy_chaincode(client, user, path, args):
    print("Deploying chaincode {} with args {}".format(path, args))
    r = client.chaincode_deploy(
        chaincode_path = path,
        type = 1,
        function = "init",
        args = args,
        secure_context = user,
    )

    s = json.dumps(r)
    print("RESPONSE " + s)

    if r["result"]["status"] != "OK":
        print("Exiting with deploy response status of ", r["status"])

    print("Deploy of chaincode {} COMPLETE".format(path))

    return r["result"]["message"]

def register_hl_user(user, key):
    print("Register hl user " + user + " with key " + key)
    data = {
      "enrollId": user,
      "enrollSecret": key
    }

    data_json = json.dumps(data)
    headers = {'Content-type': 'application/json'}
    response = requests.post("http://localhost:7050/registrar", data=data_json, headers=headers)

    print("RESPONSE : ", response)

    if response.status_code != 200:
        print("Unexpected status code in registrar " + r.status_code)
        exit(1)

    print("Response JSON " + response.text)

def stop():
    print("Tearing down poc environment...")
    teardown_docker()
    print("Docker teardown COMPLETE")


def teardown_docker():
    os.chdir("docker")
    check_call("docker-compose down", shell=True)
    check_call("docker-compose rm -f", shell=True)

commands = {
    "system_test" : system_test,
}

if __name__ == '__main__':
    command = sys.argv[1]
    if not command in commands:
        print("Unrecognized command ", command)
    else:
        commands[sys.argv[1]]()
