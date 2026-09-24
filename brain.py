import sys
import json
from NN import Brain
import torch

def getAction(floats: list[float]):
    global last_log_prob
    inputs = torch.tensor(
                floats,
                dtype=torch.float32
            )

    outputs, last_log_prob = brain.get_action(inputs)

    response = {
                "type": "action",
                "payload": outputs.tolist()
            }

    print(json.dumps(response), flush=True)

def giveReward(reward: float):
    global steps, log_probs, rewards, last_log_prob

    if last_log_prob is not None:
        log_probs.append(last_log_prob)
        rewards.append(reward)
        steps += 1

    if steps >= 1000:
        steps = 0

        # Räkna reward baklänges
        returns = []
        G = 0

        for r in reversed(rewards):
            G = r + gamma * G
            returns.append(G)

        returns.reverse()

        returns = torch.tensor(returns, dtype=torch.float32)
        returns = (
        (returns - returns.mean())
        / (returns.std() + 1e-8)
        )

        loss = -(torch.stack(log_probs) * returns).mean()
        torch.nn.utils.clip_grad_norm_(
        brain.parameters(),
        1.0
        )
        brain.optimizer.zero_grad()
        loss.backward()
        brain.optimizer.step()

        print(
            f"Training: loss={loss.item():.3f}, reward={sum(rewards):.3f}",
            file=sys.stderr
        )

        log_probs = []
        rewards = []
def train(reward:float,state:list[float]):
    giveReward(reward)
    getAction(state)

def save():
    brain.save()
brain = Brain(13, 3)
last_log_prob = None
log_probs = []
rewards = []
steps = 0
gamma = 0.99
try:
    brain.load()
except FileNotFoundError:
    brain.save()

for line in sys.stdin:
    message = json.loads(line)

    match message["type"]:
        case "getAction":
            getAction(message["payload"])
        case "giveReward":
            giveReward(message["payload"])
        case "train":
            train(message["payload"]["reward"], message["payload"]["state"])
        case "save":
            save()
            
            

            

   