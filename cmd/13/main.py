import os
from typing import NamedTuple
import numpy as np
from scipy.optimize import milp, LinearConstraint, linprog
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


def solve(A, B, P):
    c = np.array([3,1])
    b_u = np.array(P)
    b_l = np.array(P)
    A = np.array([[A[0], B[0]], [A[1], B[1]]])
    constraints = LinearConstraint(A, b_l, b_u, keep_feasible=False)
    integrality = np.ones_like(c)
    res = milp(c=c, constraints=constraints, integrality=integrality)
    if res.x is None:
        return None
    if res.x[0] > 100 or res.x[1] > 100:
        return None
    return res.fun

def part1():
    s = 0
    for claw in claws:
        c = solve(claw.A, claw.B, claw.P)
        if c is not None:
            s += c
    print(f"Part 1: {s}")


def solvep2(A, B, P):
    c = np.array([3,1])
    Ar = np.array([[A[0], B[0]], [A[1], B[1]]])
    integrality = np.ones_like(c)
    res = linprog(method='highs', c=c, A_eq=Ar, b_eq=[10000000000000 + p for p in P], integrality=integrality, options={'presolve': False})
    if res.x is None:
        return None
    return res.fun


def part2():
    s = 0
    for claw in claws:
        c = solvep2(claw.A, claw.B, claw.P)
        if c is not None:
            s += c
    print(f"Part 2: {s}")
part1()
part2()