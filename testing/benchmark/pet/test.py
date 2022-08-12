import grpc

import os
import sys
from pathlib import Path

# SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
# PARENT1 = Path(SCRIPT_DIR).parent.parent.parent.absolute().resolve() / 'api'
# PARENT = Path(SCRIPT_DIR).parent.parent.absolute().resolve()
# sys.path.append(os.path.dirname(PARENT))
# print(PARENT1)
# sys.path.append(PARENT1)
# print(os.path.dirname(PARENT1))
# for p in sys.path:
#     print(p)
# print(PARENT1)
# print(sys.path)
sys.path.append('/home/ubuntu/code/easycoding')
sys.path.append('/home/ubuntu/code/easycoding/api')
sys.path.append('/home/ubuntu/code/easycoding/api/third_party')
print(sys.path)

from api.pet import pet_pb2, rpc_grpc

def run():
    with grpc.insecure_channel('10.10.20.76:10001') as channel:
        stub = rpc_grpc.PetStoreSvcStub(channel)
        response = stub.GetPet(pet_pb2.GetPetRequest(pet_id='1'))
    print('Greeter client received: ' + response.message)


if __name__ == '__main__':
    run()
