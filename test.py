import pulp

# Define the LP problem
lp_problem = pulp.LpProblem("Optimal_Bid_Ask", pulp.LpMaximize)

# Define the decision variables
bid = pulp.LpVariable("bid", lowBound=0)
ask = pulp.LpVariable("ask", lowBound=0)

# Objective function
lp_problem += -bid + ask, "Objective"

# Add constraints
lp_problem += bid - ask <= 0
lp_problem += bid >= 26111.152617
lp_problem += ask <= 26130.947383

# Solve the problem
lp_problem.solve()

optimal_bid_value = bid.varValue
optimal_ask_value = ask.varValue
lp_status = pulp.LpStatus[lp_problem.status]

print(optimal_bid_value, optimal_ask_value, lp_status)
