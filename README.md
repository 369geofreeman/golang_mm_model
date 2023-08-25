# Inventory Control

## 1. Building the Real-Time System:

- WebSocket Connection: Establish a connection to Bybit's WebSocket API to receive real-time updates on the BTCUSDT trading pair. You'll be interested in the ticker data (for price), order book data, and recent trades.

---

## 2. Data Processing: Once we receive the data, preprocess it to extract relevant information like the mid-price, order book depth, recent trades, etc.

**Mid-Price (or Current Price):**

- Description: The price that lies exactly between the best bid and best ask in the order book.
- Usage: Used in the optimization logic to determine the bid and ask spread.
  Storage Requirement: We only need the current mid-price for our calculations, so we don't need to store a history of values.
  Recommended Data Structure: A single floating-point variable.

**Order Book Depth:**

- Description: Represents the total quantity of orders at a particular price level.
- Usage: Used to determine the base spread in our optimization function.
- Storage Requirement: We need the current order book depth for our calculations. Storing the entire order book can be memory-intensive, so it's efficient to just - compute the depth when needed.
- Recommended Data Structure: Depending on the level of granularity required, a slice or array containing the depth for each price level can be useful. For more complex operations, a balanced tree structure like a Red-Black tree or AVL tree can be effective.

**Volatility:**

- Description: Measure of how much the price of an asset varies over time.
- Usage: Used in the optimization function to determine the base spread. We're also using its exponential moving average (EMA) to determine optimization frequency.
- Storage Requirement: While we technically only need the current volatility for our base spread calculation, we're also computing an EMA, which requires us to keep - track of the previous EMA value.
- Recommended Data Structure: Two floating-point variables: one for the current volatility and another for the EMA of the volatility.

**Liquidity:**

- Description: Represents the ease with which an asset can be bought or sold without causing a significant movement in the price.
- Usage: Used in the optimization function to determine the base spread.
- Storage Requirement: We only need the current liquidity value.
- Recommended Data Structure: A single floating-point variable.

**Max Inventory:**

- Description: Represents the maximum inventory a market maker wants to hold.
- Usage: Used in the optimization function to determine bid and ask prices.
- Storage Requirement: A single value defining the limit.
- Recommended Data Structure: A single floating-point variable.
- Recent Trades (if used for volatility calculation or other metrics):

- Description: The most recent trades that have occurred in the market.
- Usage: Can be used to compute metrics like volatility if not provided directly.
- Storage Requirement: Depending on the frequency of trades and the duration over which you want to measure volatility, you might need to store a rolling window of recent trades.
- Recommended Data Structure: A circular buffer or a deque would be ideal for this, as they allow efficient additions and removals from both ends.

**Optimization Parameters:**

- Description: Parameters like alpha, beta, gamma, delta, and zeta used in the optimization function.
- Usage: Used in the optimization logic.
- Storage Requirement: Fixed set of values.
- Recommended Data Structure: Floating-point variables for each parameter.
- General Notes:

**Memory Efficiency:**

For high-frequency trading systems, memory access times can be crucial. Therefore, it's beneficial to have data structures that are cache-friendly. Arrays and slices are often better in this regard than linked structures.

**Concurrency:**

If multiple parts of the system access and modify the data simultaneously, consider using synchronization mechanisms or concurrent data structures to ensure data consistency.

Data Freshness: In fast-moving markets, data can become stale quickly. It's essential to have a mechanism to update data rapidly and ensure the market maker operates on the most recent data.

---

## 3. Optimization Logic: Using the real-time data, run your optimization logic to determine the optimal bid-ask spread based on the market-making model we discussed.

- Parameter Management:
  The parameters (alpha, beta, gamma, delta, zeta) are initialized and can be dynamically adjusted. This gives flexibility to the strategy.
- Validation:
  validateParameters checks if the input parameters are within logical ranges, ensuring no anomalies or corrupted data influence the optimization process.
- EMA for Volatility:
  The volatility is being tracked using an Exponential Moving Average (EMA) which provides a smoothed representation of volatility. This aids in determining the frequency of optimization.
- Optimization Logic:
  Linear programming is used to determine the optimal bid and ask spread based on given constraints and an objective function.
  In case the optimization fails, the function reverts to previously successful values or a default spread.
- Optimization Frequency:
  getOptimizationFrequency determines how often the optimization should run based on the EMA of the volatility.
- Cost Function & Base Spread:
  The costFunction calculates the risk associated with the inventory and deviation from the current price.
  The baseSpreadFunction computes the base spread considering volatility, liquidity, and order book depth.

- Display/Logging: For monitoring, display the optimal bid and ask prices, the current market price, and other relevant metrics in real-time. Additionally, log this data for future analysis.

- Risk Management: Even though you won't be executing trades, it's still essential to simulate risk management practices. For instance, track a hypothetical inventory and see how it evolves over time.

---

## 4. Building the Backtesting System:

- Data Loading: Load the historical data from Bybit. This data should be at the same granularity as the real-time data you'll be using (e.g., tick-by-tick or minute-by-minute).

- Simulation Loop: Loop through the historical data, simulating the passage of time. At each time step, use the data to determine the optimal bid-ask spread, just as in the real-time system.

- Trade Execution Simulation: If your bid or ask gets "hit" (i.e., if the market price reaches your bid or ask), simulate the execution of a trade. Adjust your hypothetical inventory accordingly.

- Metrics Calculation: At the end of the backtest, calculate performance metrics like total profit, Sharpe ratio, maximum drawdown, etc.

- Parameter Tuning: Use the backtesting system to optimize the parameters of your market-making model. This might involve grid search, gradient-based optimization, or other methods.

---

## 5. Other notes

Alpha: Influences the inventory risk component.
Beta: Influences the price risk component.
Gamma: Affects the volatility term in the base spread function.
Delta: Affects the liquidity term in the base spread function.
Zeta: Affects the order book depth term in the base spread function.

α: Represents the weight associated with inventory risk in the cost function. A higher
α means that the model will place more emphasis on minimizing inventory risk when determining the optimal spread.

β: Represents the weight associated with price risk in the cost function. A higher
β indicates that deviations from the current price (when setting bid and ask prices) are considered more risky and costly.

γ: Determines the impact of market volatility on the base spread. A higher
γ means the model will increase the spread more in response to increased volatility.

δ: Represents the influence of liquidity on the base spread. If
δ is high, then lower liquidity will lead to a wider spread.

ζ: Determines the significance of order book depth on the base spread. A higher
ζ means that the model considers a deeper order book as having a more significant impact on the spread.

Data Source: Ensure that you have access to high-quality historical data, including order book snapshots, trades, volatility, and liquidity metrics. The granularity of the data (e.g., tick-by-tick, minute-by-minute) will influence the accuracy and applicability of the model.

Simulation Environment: It's important to have a realistic backtesting environment that can simulate order execution, slippage, and other real-world trading frictions.

Adverse Selection: Market makers face the risk of adverse selection, where traders with better information trade against them. This is an inherent risk of market making and can be managed but not eliminated.

Parameters Initialization: While the optimization process will fine-tune the parameters, having a good starting point or a reasonable range for the grid search can speed up the process.

Evaluation Metrics: Besides profitability, consider other metrics like the Sharpe ratio, maximum drawdown, and turnover rate to evaluate the strategy's performance.

Robustness Checks: Once you identify a set of parameters that work well on your historical data, it's essential to test the strategy on out-of-sample data or use techniques like bootstrapping to ensure the strategy's robustness.

Market Regime Changes: Cryptocurrency markets can undergo regime changes, where the market's behavior shifts due to macroeconomic factors, regulatory changes, or other reasons. It's beneficial to segment the data into different regimes and see how the strategy performs in each.

Liquidity Concerns: Ensure that the strategy doesn't assume infinite liquidity. In reality, large orders can move the market.

Operational Risks: While not directly related to the strategy, it's crucial to understand the operational risks associated with real-time trading, like system outages, exchange downtimes, and more.

Feedback Effects: A large market maker can influence the market. Ensure that the strategy accounts for its potential impact on the market.

---

# Misc

---

set of values.

- Recommended Data Structure: Floating-point variables for each parameter.

The market maker will use an optimization function to handle most of the logic and have a cost function associated with it. Here are the details of the optimization logic:

- Parameter Management:
  The parameters (alpha, beta, gamma, delta, zeta) are initialized and can be dynamically adjusted. This gives flexibility to the strategy.
- Validation:
  validateParameters checks if the input parameters are within logical ranges, ensuring no anomalies or corrupted data influence the optimization process.
- EMA for Volatility:
  The volatility is being tracked using an Exponential Moving Average (EMA) which provides a smoothed representation of volatility. This aids in determining the frequency of optimization.
- Optimization Logic:
  Linear programming is used to determine the optimal bid and ask spread based on given constraints and an objective function.
  In case the optimization fails, the function reverts to previously successful values or a default spread.
- Optimization Frequency:
  getOptimizationFrequency determines how often the optimization should run based on the EMA of the volatility.
- Cost Function & Base Spread:
  The costFunction calculates the risk associated with the inventory and deviation from the current price.
  The baseSpreadFunction computes the base spread considering volatility, liquidity, and order book depth.

- Display/Logging: For monitoring, display the optimal bid and ask prices, the current market price, and other relevant metrics in real-time. Additionally, log this data for future analysis.

- Risk Management: Even though we won't be executing trades, it's still essential to simulate risk management practices. For instance, track a hypothetical inventory and see how it evolves over time.

The parameters (alpha, beta, gamma, delta, zeta) are defined below:
α: Represents the weight associated with inventory risk in the cost function. A higher
α means that the model will place more emphasis on minimizing inventory risk when determining the optimal spread.

β: Represents the weight associated with price risk in the cost function. A higher
β indicates that deviations from the current price (when setting bid and ask prices) are considered more risky and costly.

γ: Determines the impact of market volatility on the base spread. A higher
γ means the model will increase the spread more in response to increased volatility.

δ: Represents the influence of liquidity on the base spread. If
δ is high, then lower liquidity will lead to a wider spread.

ζ: Determines the significance of order book depth on the base spread. A higher
ζ means that the model considers a deeper order book as having a more significant impact on the spread.
