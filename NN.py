import torch
import torch.nn as nn
import torch.optim as optim



class Brain(nn.Module):

    def __init__(self, inputs = 15, outputs=2):
        super().__init__()
        self.network = nn.Sequential(
        nn.Linear(inputs, 32),
        nn.ReLU(),
        nn.Linear(32, 32),
        nn.ReLU(),
        nn.Linear(32, outputs),
        nn.Tanh()
        )
        self.inputs= inputs
        self.outputs = outputs
        self.optimizer = optim.Adam(self.network.parameters(), lr=0.003)


    def save(self):
        # name it the with how many inputs and outputs it has
        torch.save(self.state_dict(), f"brain-{self.inputs}x{self.outputs}.pth")  
        # print("fil sparad")

    def load(self):
        state = torch.load(
        f"brain-{self.inputs}x{self.outputs}.pth",
        weights_only=True
        )
        self.load_state_dict(state)

    def get_action(self, state):
        mean = self.forward(state)

        dist = torch.distributions.Normal(mean, 0.2)

        action = dist.sample()

        log_prob = dist.log_prob(action).sum()

        return action, log_prob

    def forward(self, state : nn.Sequential): 
        return self.network(state)