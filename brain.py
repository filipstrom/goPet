import sys
import json
from NN import Brain
import torch

brain = Brain(13,2)
for line in sys.stdin:


    inputs = torch.tensor(json.loads(line) , dtype=torch.float32)

    outputs = brain(inputs)

    print(json.dumps(outputs.tolist()), flush=True)