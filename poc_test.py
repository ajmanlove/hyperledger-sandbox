import os
import os.path
# from subprocess import check_call
import sys
# from hyperledger.client import Client
import requests
import json

import time

import configparser

setup_hl_creds = ("test_user0", "MS9qrN8hFjlE")
insurer1_hl_creds = ("insurer1", "iF6bJY5sWumg")
reinsurer1_hl_creds = ("reinsurer1", "8jUZV34kUBcN")
reinsurer2_hl_creds = ("reinsurer2", "RmMio3LCk37J")
insurer2_hl_creds = ("insurer2", "1wtM0CXjIVXr")
reinsurer3_hl_creds = ("reinsurer3", "atjQRL2S6FJx")

cch = "CHAINCODE"
am = "asset_management"
rr = "reinsurance_request"
rp = "reinsurance_proposal"

def system_test():
    if not os.path.exists(".setup.ini"):
        print("poc_test requires the .setup.ini file written by setup_poc.py")
        exit(1)

    config = configparser.ConfigParser()
    config.read('.setup.ini')
    print(config.sections())

    # c = Client(base_url="http://127.0.0.1:7050")

    subId = submit(config)

    ri1prop = propose(config, "reinsurer1", subId, "insurer1")

    ## Ensure reinsurer2 (the other requestee) cannot view the proposal
    r = get_proposal(config, "reinsurer2", ri1prop, post)
    assert 'error' in r
    assert 'Insuffienct rights on asset' in r['error']['data']

    ri2prop = propose(config, "reinsurer2", subId, "insurer1")
    ## Ensure reinsurer1 (the other requestee) cannot view the proposal
    r = get_proposal(config, "reinsurer1", ri2prop, post)
    assert 'error' in r
    assert 'Insuffienct rights on asset' in r['error']['data']

    ## Ensure insurer1 can counter
    counter(config, "insurer1", ri1prop, "reinsurer1")

    ## Ensure reinsurer1 can counter
    counter(config, "reinsurer1", ri1prop, "insurer1")

    ## Ensure reinsurer1 can't accept or reject
    ## TODO invoke does not return the error condition in REST
    # response = post_reject(config, "reinsurer1", ri2prop, post)
    # reject(config, "reinsurer1", ri2prop, "insurer1")

    ## Reject reinsurer2 proposal
    reject(config, "insurer1", ri2prop, "reinsurer2")

    ## Accept reinsure1 proposal
    accept(config, "insurer1", ri1prop, "reinsurer1")

    print("System Test COMPLETE")
    exit(1)

## submit request, return submission id
def submit(config):
    print("submit()")
    data = {
        "jsonrpc": "2.0",
        "method": "invoke",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config[cch][rr]
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

    ## Invoke the submission
    assert_post(data)
    time.sleep(2) ## TODO try to avoid a sleep

    ## Get insurer1 user assets and verify
    i1ua = to_json(get_user_assets(config, "insurer1"))
    ## Get the last submission id
    keys = i1ua['submissions'].keys()
    last = sorted(keys, key=subm_sort_key)[len(keys) - 1]

    ## Get insurer1 asset rights and verify
    i1ar = to_json(get_asset_rights(config, 'insurer1', last))
    assert i1ar == {'Rights': [0, 1], 'Exists': True}

    ## Get and verify reinsurer 1 asset rights and user assets
    ri1ar = to_json(get_asset_rights(config, 'reinsurer1', last))
    assert ri1ar == {'Rights': [1], 'Exists': True}
    ri1ua = to_json(get_user_assets(config, "reinsurer1"))
    assert last in ri1ua['requests']


    ## Get the submission record and verify
    submission = to_json(get_submission(config, "insurer1", last))
    assert submission["id"] == last
    ## TODO verify whole submission object

    ## Ensure reinsurer1 and reinsurer2 can view the submission
    assert submission == to_json(get_submission(config, "reinsurer1", last))
    assert submission == to_json(get_submission(config, "reinsurer2", last))

    ## Ensure reinsurer3 and insurer2 cannot view the submission
    r = get_submission(config, "reinsurer3", last, post)
    assert 'error' in r
    assert 'Insuffienct rights on asset' in r['error']['data']

    r = get_submission(config, "insurer2", last, post)
    assert 'error' in r
    assert 'Insuffienct rights on asset' in r['error']['data']

    return last


def get_submission(config, user, id, poster=None):
    print("get_submission() [{0}, {1}]".format(user, id))
    data = {
        "jsonrpc": "2.0",
        "method": "query",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config[cch][rr]
            },
            "ctorMsg": {
                "function": "get_request",
                "args": [
                    id
                ]
            },
            "secureContext": user,
            "attributes": ["enrollmentId"]
        },
        "id": 3
    }

    if poster is None:
        return assert_post(data)
    else:
        return poster(data)

def subm_sort_key(v):
    return int(v.split("-")[1])

def propose(config, user, requestId, requestUser):
    print("propose()")
    data = {
        "jsonrpc": "2.0",
        "method": "invoke",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config[cch][rp]
            },
            "ctorMsg": {
                "function": "propose",
                "args": [
                    requestId, "some user {0} text".format(user)
                ]
            },
            "secureContext": user,
            "attributes": ["enrollmentId"]
        },
        "id": 2
    }

    ## Propose
    assert_post(data)
    time.sleep(2) ## TODO try to avoid a sleep

    ## Get insurer1 user assets and verify
    ## TODO verify
    proposerua = to_json(get_user_assets(config, user))
    propId = None
    for k,v in proposerua['proposals'].items():
        if v['submissionId'] == requestId:
            propId = k

    assert propId is not None

    ## Get the proposal and verify
    proposal = to_json(get_proposal(config, user, propId))
    assert propId == proposal['id']
    assert requestId == proposal['requestId']
    assert user == proposal['bidder']
    assert user == proposal['updatedBy']
    assert "bid" == proposal['status']

    ## Assert the requesting user can access the proposal
    assert proposal == to_json(get_proposal(config, requestUser, propId))

    ## verify that user assets have changed (updated, updatedBy)
    ua = to_json(get_user_assets(config, user))
    assert propId in ua['proposals']
    assert requestId == ua['proposals'][propId]['submissionId']
    assert propId == ua['proposals'][propId]['proposalId']
    assert user in ua['proposals'][propId]['updatedBy']

    ua = to_json(get_user_assets(config, requestUser))
    assert propId in ua['proposals']
    assert requestId == ua['proposals'][propId]['submissionId']
    assert propId == ua['proposals'][propId]['proposalId']
    assert user in ua['proposals'][propId]['updatedBy']

    return propId

def get_proposal(config, user, id, poster=None):
    print("get_proposal() [{0}, {1}]".format(user, id))

    data = {
        "jsonrpc": "2.0",
        "method": "query",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config[cch][rp]
            },
            "ctorMsg": {
                "function": "get_proposal",
                "args": [
                    id
                ]
            },
            "secureContext": user,
            "attributes": ["enrollmentId"]
        },
        "id": 3
    }

    if poster is None:
        return assert_post(data)
    else:
        return poster(data)

def counter(config, user, propId, userB):
    print("counter()")
    counterText = "COUNTER text by user {0}".format(user)
    data = {
        "jsonrpc": "2.0",
        "method": "invoke",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config[cch][rp]
            },
            "ctorMsg": {
                "function": "counter",
                "args": [
                    propId, counterText
                ]
            },
            "secureContext": user,
            "attributes": ["enrollmentId"]
        },
        "id": 2
    }

    ## Propose
    assert_post(data)
    time.sleep(2) ## TODO try to avoid a sleep

    ## Get the proposal and verify
    proposal = to_json(get_proposal(config, user, propId))
    assert user == proposal['updatedBy']
    assert counterText == proposal['contractText']
    assert "counter" == proposal['status']

    ## Assert the requesting user can access the proposal
    assert proposal == to_json(get_proposal(config, userB, propId))

    ## verify that user assets have changed (updated, updatedBy)
    ua = to_json(get_user_assets(config, user))
    assert propId in ua['proposals']
    assert user in ua['proposals'][propId]['updatedBy']

    ua = to_json(get_user_assets(config, userB))
    assert propId in ua['proposals']
    assert user in ua['proposals'][propId]['updatedBy']


def accept(config, user, propId, acceptedUser):
    print("accept()")

    data = {
        "jsonrpc": "2.0",
        "method": "invoke",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config[cch][rp]
            },
            "ctorMsg": {
                "function": "accept",
                "args": [
                    propId
                ]
            },
            "secureContext": user,
            "attributes": ["enrollmentId"]
        },
        "id": 2
    }

    ##post
    assert_post(data)
    time.sleep(2) ## TODO try to avoid a sleep

    ## Get the proposal and verify
    proposal = to_json(get_proposal(config, user, propId))
    assert user == proposal['updatedBy']
    assert "accepted" == proposal['status']

    ## Assert the rejected user can access the proposal
    assert proposal == to_json(get_proposal(config, acceptedUser, propId))

    ## verify user assets
    ua = to_json(get_user_assets(config, acceptedUser))
    assert propId not in ua['proposals']
    assert propId in ua['accepted']

    ua = to_json(get_user_assets(config, user))
    assert propId not in ua['proposals']
    assert propId in ua['accepted']

def post_reject(config, user, propId, poster):
    print("reject()")

    data = {
        "jsonrpc": "2.0",
        "method": "invoke",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config[cch][rp]
            },
            "ctorMsg": {
                "function": "reject",
                "args": [
                    propId
                ]
            },
            "secureContext": user,
            "attributes": ["enrollmentId"]
        },
        "id": 2
    }

    ##post
    return poster(data)


def reject(config, user, propId, rejectedUser):
    print("reject()")
    post_reject(config, user, propId, assert_post)
    time.sleep(2) ## TODO try to avoid a sleep

    ## Get the proposal and verify
    proposal = to_json(get_proposal(config, user, propId))
    assert user == proposal['updatedBy']
    assert "rejected" == proposal['status']

    ## Assert the rejected user can access the proposal
    assert proposal == to_json(get_proposal(config, rejectedUser, propId))

    ## verify user assets
    ua = to_json(get_user_assets(config, rejectedUser))
    assert propId not in ua['proposals']
    assert propId in ua['rejected']

    ua = to_json(get_user_assets(config, user))
    assert propId not in ua['proposals']
    assert propId in ua['rejected']



def get_user_assets(config, user):
    print("get_user_assets() [{0}]".format(user))
    data = {
        "jsonrpc": "2.0",
        "method": "query",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config[cch][am]
            },
            "ctorMsg": {
                "function": "get_user_assets",
                "args": [
                ]
            },
            "secureContext": user,
            "attributes": ["enrollmentId"]
        },
        "id": 3
    }

    return assert_post(data)

def get_asset_rights(config, user, asset):
    print("get_asset_rights() [{0}, {1}]".format(user, asset))
    data = {
        "jsonrpc": "2.0",
        "method": "query",
        "params": {
            "type": 1,
            "chaincodeID": {
                "name": config['CHAINCODE'][am]
            },
            "ctorMsg": {
                "function": "get_asset_rights",
                "args": [
                    user, asset
                ]
            },
            "secureContext": user,
            "attributes": ["enrollmentId"]
        },
        "id": 3
    }

    return assert_post(data)

def post(data):
    data_json = json.dumps(data)
    headers = {'Content-type': 'application/json'}
    response = requests.post("http://localhost:7050/chaincode", data=data_json, headers=headers)

    if response.status_code != 200:
        print("Unexpected status code in registrar " + response.status_code)
        exit(1)

    # print("Response JSON " + response.text)

    return response.json()

def assert_post(data):
    res = post(data)
    assert res['result']['status'] == 'OK'
    return res['result']['message']

def to_json(string):
    return json.loads(string)

commands = {
    "system_test" : system_test,
}

if __name__ == '__main__':
    command = sys.argv[1]
    if not command in commands:
        print("Unrecognized command ", command)
    else:
        commands[sys.argv[1]]()
