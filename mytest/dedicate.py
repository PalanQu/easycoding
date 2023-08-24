import json
import requests
import time


# with open('dedicate.json') as f:
#     dedicate_json = json.load(f)
#     for res in dedicate_json['data']['result']:
#         if 'label_dedicate' in res['metric']:
#             print(res['metric']['label_dedicate'])

BIZ_BLACK_LIST = [
    'agent-smith',
    'red-event-logger',
    'hdfs-tmp',
]

def remove_black_list_items(array1):
    return [item for item in array1 if item not in BIZ_BLACK_LIST]

with open('./dedicates.txt') as f:
    for dedicate in f.read().splitlines():
        PQL1 = 'node_filesystem_avail_bytes {mountpoint="/data/agentsmith", xhs_zone!="qcsh5-qa", env!="staging",label_env!="staging",address_InternalIP!="10.214.21.134", label_dedicate="' + dedicate + '"}'
        params1 = {
            'query': PQL1,
            'time': '1692857559'
        }
        headers1 = {
            'authority': 'monitor.devops.xiaohongshu.com',
            'cookie': 'access-token-oa.xiaohongshu.com=internal.xhsoa.AT-3976d145b33a46f783a8e0a9aa71cbe2-0c54ecb1e8bd4921a90aab5cf7b03292; redpass_did=E6E26190-1939-5D23-BA62-B7C5E37EF323; xhsTrackerId=3a525a7a-50c3-41a9-ba72-cdaa1f69c114; xhsTrackerId.sig=TEoKpRBawqfgXa6aNkEz1N4AlyROxiy5GitouCyCwJ8; xsecappid=xhs-pc-web; a1=186b5033c262dq5c3fntdr7h2b3fb7o03xmvxmoqn30000212451; webId=40848649ee2761868f1f35c644c3def3; gid=yYKD28qqiKqdyYKD28qqSC6lJKJfA2SqiI9f6AIFWxJDqhq8iDWl8C888JyJ42y8Y0WiDfiW; gid.sign=MVhNrZ8tYXvHObwRQEAedqCU+s0=; web_session=030037a336356249065b5a277d234a8ca2b512; _ga=GA1.2.941828200.1680767657; x-user-id=63e20cbd65d5fa0001e131ef; access-token-city.sit.xiaohongshu.com=internal.xhsoa.AT-9354d8f6ed3247b5ba2ba6e3fd11f22a-ac3642163d804df6a2a36e62fc2e7090; access-token-redcloud.devops.xiaohongshu.com=internal.tech.AT-b150dadbada040769e1d2a1568560f98-514dbce0b3b34ac48a4d7c9db283d83c; access-token-redserverless.devops.sit.xiaohongshu.com=internal.redserverless.AT-fdfccbe2f7774ceda5503a05b4eb343c-9d8b44483b2544baaf9907fb368704fb; porch_beaker_session_id=bf65b8d145709b4c46904beb295776dfb6d33ca4gAJ9cQAoWAMAAABfaWRxAVggAAAAMGVjZjgyZmU2MTBkNGY5MjhkNzkwNGFlZWE2NzhlNGJxAlgOAAAAX2FjY2Vzc2VkX3RpbWVxA0rPGZZkWAoAAABleHBpcmVkX2F0cQRKz8AMZVgNAAAAcG9yY2gtdXNlci1pZHEFWBgAAAA2M2UyMGNiZDY1ZDVmYTAwMDFlMTMxZWZxBlgQAAAAcG9yY2gtYXV0aC10b2tlbnEHWEEAAAA0MWM5MTUyMjdjNTA0ZjQ0YWU4YzJlMWExNzcyNTE2ZC1iNjg0ZTA0OTJjMTc0OWI2YjdhMjEyZDdhNjY1MWZlZHEIWA4AAABfY3JlYXRpb25fdGltZXEJSs8ZlmR1Lg==; _ga_T11SF3WXX2=GS1.2.1687687393.1.0.1687687393.60.0.0; _ga_K2SPJK2C73=GS1.2.1687687394.1.0.1687687394.60.0.0; access-token-its-m.xiaohongshu.com=internal.bail.AT-2b87d4e53f164228b1f0f955f069f2c4-3b43abc0808c4096999a450392319fe9; access-token-yunxiao.devops.xiaohongshu.com=internal.yanxiao.AT-0eefbaf11b8c49c09b1119b88fda963d-6ba1fcff20bf46708b4b78a64dd23f24; access-token-hr.xiaohongshu.com=internal.xhsoa.AT-52fbf92951dc4fb5952ad3190cf3d212-eeb59938aca648de91ff5920f672d08f; access-token-its.xiaohongshu.com=internal.bail.AT-bc3388131ad54713b0920a9d393c39ff-08bf8c44617342a2a974db1f0ef61fa7; access-token-city.xiaohongshu.com=internal.xhsoa.AT-47e18444f891478eb4101dad083dbf32-e641a700a1bf4a4f9e8a0530e38ddaa4; sit_porch_beaker_session_id=5acffe1d861dd1a6b27dd379edb1dc3546813045KGRwMQpWcG9yY2gtdXNlci1pZApwMgpWNjNlMjBjYmQ2NWQ1ZmEwMDAxZTEzMWVmCnAzCnNWX2lkCnA0ClY0YTY3YmUwMDNjYTQxMWVlYjA0MWVmYjkwNWI1ZGYxOQpwNQpzVl9hY2Nlc3NlZF90aW1lCnA2CkYxNjkyMjM4NjgxLjU2ODAwMDEKc1Zwb3JjaC1hdXRoLXRva2VuCnA3ClZhZGZhOTUzNGE1ZWI0ZWJhYTExYmE0YTNkOThhMjk5Ny02ZjYwY2RmOWYyNDY0Mzc1OTBjMjU3OTEwMzAxMjdhMApwOApzVl9jcmVhdGlvbl90aW1lCnA5CkYxNjkyMjM4NjgxLjU2ODAwMDEKc1ZleHBpcmVkX2F0CnAxMApJMTcwMDAxNDY4MQpzLg==; timeslogin=98108d204856ba012bc918f3ad137733; newresults=; _porch_uuid=732026864162461; abRequestId=40848649ee2761868f1f35c644c3def3; webBuild=3.5.1; xhs_sre_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoicXVqaWFiYW9AeGlhb2hvbmdzaHUuY29tIiwibmFtZSI6IuWkqeepuSjmm7LlmInlrp0pIiwiZW1haWwiOiJ0aWFucWlvbmdAeGlhb2hvbmdzaHUuY29tIiwidW0iOiJ0aWFucWlvbmciLCJsb2dvIjoiaHR0cHM6Ly93ZXdvcmsucXBpYy5jbi93d3BpYy8zNzYwMjVfVlhIVHU4dEdTeGVGcVZ2XzE2NjM5ODQyNzUvMCIsImRlcGFydG1lbnRfaWQiOjQ5MzYsInZlcnNpb24iOiIyMDIzLTAzLTAxIiwiaXNzIjoiemhpYmFuIiwiZXhwIjoxNjkyODYwODg5LCJpYXQiOjE2OTI2MDE2ODl9.0K-kPl_Ii-opU9kfxJ565ttYeHoL-MtYhixadGLfQoA; grafana_session=11cc08f4880d095d16856ca68e23498f; hulk.sso.jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoVG9rZW4iOiI0MWM5MTUyMjdjNTA0ZjQ0YWU4YzJlMWExNzcyNTE2ZC1iNjg0ZTA0OTJjMTc0OWI2YjdhMjEyZDdhNjY1MWZlZCIsInVzZXJJZCI6IjYzZTIwY2JkNjVkNWZhMDAwMWUxMzFlZiIsInVzZXJTdWJzeXN0ZW1zSW5mbyI6e30sImlhdCI6MTY5MjY4ODAxMCwiZXhwIjoxNjkzMjkyODEwfQ.lotA4QnbZsGNGem8nsFMYhaVbmZ-BNNlGSu7lYnTrFU; R_SESS=token-hthcg:kmp7m5lrrxtxhz92bt2jfjgqwstmz2qzc2dxfzvv7qqpbxlf9wbrx4',
        }
        url1 = 'https://monitor.devops.xiaohongshu.com/api/datasources/proxy/32/api/v1/query'
        response1 = requests.get(url1, headers=headers1, params=params1)
        no_biz = True
        for res in response1.json()['data']['result']:
            node_ip = res['metric']['address_InternalIP']
            PQL2 = 'count(agentsmith_value{nodeIP="' + node_ip  +'"}) by (biz)'
            params2 = {
                'query': PQL2,
                'time': '1692857559'
            }
            headers2 = {
                'authority': 'monitor.devops.xiaohongshu.com',
                'datasource-name': 'vms-agentsmith',
                'cookie': 'access-token-oa.xiaohongshu.com=internal.xhsoa.AT-3976d145b33a46f783a8e0a9aa71cbe2-0c54ecb1e8bd4921a90aab5cf7b03292; redpass_did=E6E26190-1939-5D23-BA62-B7C5E37EF323; xhsTrackerId=3a525a7a-50c3-41a9-ba72-cdaa1f69c114; xhsTrackerId.sig=TEoKpRBawqfgXa6aNkEz1N4AlyROxiy5GitouCyCwJ8; xsecappid=xhs-pc-web; a1=186b5033c262dq5c3fntdr7h2b3fb7o03xmvxmoqn30000212451; webId=40848649ee2761868f1f35c644c3def3; gid=yYKD28qqiKqdyYKD28qqSC6lJKJfA2SqiI9f6AIFWxJDqhq8iDWl8C888JyJ42y8Y0WiDfiW; gid.sign=MVhNrZ8tYXvHObwRQEAedqCU+s0=; web_session=030037a336356249065b5a277d234a8ca2b512; _ga=GA1.2.941828200.1680767657; x-user-id=63e20cbd65d5fa0001e131ef; access-token-city.sit.xiaohongshu.com=internal.xhsoa.AT-9354d8f6ed3247b5ba2ba6e3fd11f22a-ac3642163d804df6a2a36e62fc2e7090; access-token-redcloud.devops.xiaohongshu.com=internal.tech.AT-b150dadbada040769e1d2a1568560f98-514dbce0b3b34ac48a4d7c9db283d83c; access-token-redserverless.devops.sit.xiaohongshu.com=internal.redserverless.AT-fdfccbe2f7774ceda5503a05b4eb343c-9d8b44483b2544baaf9907fb368704fb; porch_beaker_session_id=bf65b8d145709b4c46904beb295776dfb6d33ca4gAJ9cQAoWAMAAABfaWRxAVggAAAAMGVjZjgyZmU2MTBkNGY5MjhkNzkwNGFlZWE2NzhlNGJxAlgOAAAAX2FjY2Vzc2VkX3RpbWVxA0rPGZZkWAoAAABleHBpcmVkX2F0cQRKz8AMZVgNAAAAcG9yY2gtdXNlci1pZHEFWBgAAAA2M2UyMGNiZDY1ZDVmYTAwMDFlMTMxZWZxBlgQAAAAcG9yY2gtYXV0aC10b2tlbnEHWEEAAAA0MWM5MTUyMjdjNTA0ZjQ0YWU4YzJlMWExNzcyNTE2ZC1iNjg0ZTA0OTJjMTc0OWI2YjdhMjEyZDdhNjY1MWZlZHEIWA4AAABfY3JlYXRpb25fdGltZXEJSs8ZlmR1Lg==; _ga_T11SF3WXX2=GS1.2.1687687393.1.0.1687687393.60.0.0; _ga_K2SPJK2C73=GS1.2.1687687394.1.0.1687687394.60.0.0; access-token-its-m.xiaohongshu.com=internal.bail.AT-2b87d4e53f164228b1f0f955f069f2c4-3b43abc0808c4096999a450392319fe9; access-token-yunxiao.devops.xiaohongshu.com=internal.yanxiao.AT-0eefbaf11b8c49c09b1119b88fda963d-6ba1fcff20bf46708b4b78a64dd23f24; access-token-hr.xiaohongshu.com=internal.xhsoa.AT-52fbf92951dc4fb5952ad3190cf3d212-eeb59938aca648de91ff5920f672d08f; access-token-its.xiaohongshu.com=internal.bail.AT-bc3388131ad54713b0920a9d393c39ff-08bf8c44617342a2a974db1f0ef61fa7; access-token-city.xiaohongshu.com=internal.xhsoa.AT-47e18444f891478eb4101dad083dbf32-e641a700a1bf4a4f9e8a0530e38ddaa4; sit_porch_beaker_session_id=5acffe1d861dd1a6b27dd379edb1dc3546813045KGRwMQpWcG9yY2gtdXNlci1pZApwMgpWNjNlMjBjYmQ2NWQ1ZmEwMDAxZTEzMWVmCnAzCnNWX2lkCnA0ClY0YTY3YmUwMDNjYTQxMWVlYjA0MWVmYjkwNWI1ZGYxOQpwNQpzVl9hY2Nlc3NlZF90aW1lCnA2CkYxNjkyMjM4NjgxLjU2ODAwMDEKc1Zwb3JjaC1hdXRoLXRva2VuCnA3ClZhZGZhOTUzNGE1ZWI0ZWJhYTExYmE0YTNkOThhMjk5Ny02ZjYwY2RmOWYyNDY0Mzc1OTBjMjU3OTEwMzAxMjdhMApwOApzVl9jcmVhdGlvbl90aW1lCnA5CkYxNjkyMjM4NjgxLjU2ODAwMDEKc1ZleHBpcmVkX2F0CnAxMApJMTcwMDAxNDY4MQpzLg==; timeslogin=98108d204856ba012bc918f3ad137733; newresults=; _porch_uuid=732026864162461; abRequestId=40848649ee2761868f1f35c644c3def3; webBuild=3.5.1; grafana_session=11cc08f4880d095d16856ca68e23498f; hulk.sso.jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoVG9rZW4iOiI0MWM5MTUyMjdjNTA0ZjQ0YWU4YzJlMWExNzcyNTE2ZC1iNjg0ZTA0OTJjMTc0OWI2YjdhMjEyZDdhNjY1MWZlZCIsInVzZXJJZCI6IjYzZTIwY2JkNjVkNWZhMDAwMWUxMzFlZiIsInVzZXJTdWJzeXN0ZW1zSW5mbyI6e30sImlhdCI6MTY5MjY4ODAxMCwiZXhwIjoxNjkzMjkyODEwfQ.lotA4QnbZsGNGem8nsFMYhaVbmZ-BNNlGSu7lYnTrFU; R_SESS=token-hthcg:kmp7m5lrrxtxhz92bt2jfjgqwstmz2qzc2dxfzvv7qqpbxlf9wbrx4',
                'referer': 'https://monitor.devops.xiaohongshu.com/explore?orgId=13&left=%5B%22now-30m%22,%22now%22,%22vms-agentsmith%22,%7B%22exemplar%22:true,%22expr%22:%22count(agentsmith_value%7BnodeIP%3D%5C%2210.11.32.169%5C%22%7D)%20by%20(biz)%22,%22requestId%22:%22Q-3acd9ddd-811b-48a3-bafd-9ecd848a5073-0A%22%7D%5D'
            }
            url2 = 'https://monitor.devops.xiaohongshu.com/api/datasources/proxy/135/api/v1/query'
            response2 = requests.get(url2, headers=headers2, params=params2)
            biz_list = []
            for biz in response2.json()['data']['result']:
                if 'biz' in biz['metric']:
                    biz_name = biz['metric']['biz']
                    biz_list.append(biz_name)

            biz_list = remove_black_list_items(biz_list)
            print('dedicate: {}     node_ip: {}     biz_list: {}'.format(dedicate, node_ip, biz_list))
            time.sleep(1)
            if len(biz_list) > 0:
                no_biz = False
                with open('./tags.txt', 'a') as f1:
                    f1.write(dedicate)
                    f1.write('\n')
                break
        if no_biz:
            with open('./tags2.txt', 'a') as f1:
                f1.write(dedicate)
                f1.write('\n')


