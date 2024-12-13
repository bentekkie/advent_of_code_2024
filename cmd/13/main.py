import os
from typing import NamedTuple
import numpy as np
from scipy.optimize import milp, LinearConstraint, linprog, Bounds
from argparse import ArgumentParser

'''
Button A: X+94, Y+34
Button B: X+22, Y+67
Prize: X=8400, Y=5400
'''

THIS_FILE = os.path.abspath(__file__)
THIS_DIR = os.path.dirname(os.path.abspath(__file__))


class Claw(NamedTuple):
    A: list[int]
    B: list[int]
    P: list[int]


claws : list[Claw] = []

parser = ArgumentParser()
parser.add_argument("--use_example", action='store_true')
args = parser.parse_args()

fname = 'example.txt' if args.use_example else  'input.txt'

day = os.path.basename(THIS_DIR)

with open(os.path.join(THIS_DIR, '..','..','pkg','inputs','puzzle_inputs',day,fname)) as f:
    A = []
    B = []
    P = []
    for line in f.readlines():
        if line.startswith("Button A: "):
            A = [int(x) for x in line.removeprefix("Button A: ").replace("X+", "").replace("Y+", "").split(", ")]
        elif line.startswith("Button B: "):
            B = [int(x) for x in line.removeprefix("Button B: ").replace("X+", "").replace("Y+", "").split(", ")]
        elif line.startswith("Prize: "):
            P = [int(x) for x in line.removeprefix("Prize: ").replace("X=", "").replace("Y=", "").split(", ")]
        else:
            claws.append(Claw(A,B,P))
            A = []
            B = []
            P = []
    claws.append(Claw(A,B,P))


def solve(A, B, P, bounds=None):
    res = linprog(
        method='highs',
        c=[3,1],
        A_eq=list(zip(A,B)),
        b_eq=P,
        bounds=bounds, 
        integrality=[1, 1],
        options={'presolve': False},
    )
    if res.x is None:
        return None
    return res.fun

def part1():
    s = 0
    for claw in claws:
        c = solve(claw.A, claw.B, claw.P, bounds=(0, 100))
        if c is not None:
            s += c
    print(f"Part 1: {s}")

def part2():
    s = 0
    for claw in claws:
        c = solve(claw.A, claw.B, [10000000000000 + p for p in claw.P])
        if c is not None:
            s += c
    print(f"Part 2: {s}")
part1()
part2()