# Exploring Concurrency
This assignment asked students to explore using Go's concurrency features (channels and goroutines in particular) on a program that fits all possible linear regretions for the Boston Housing Study. The code considers only main-effect models, with no-interactions and no transformations of the explanatory variables. The models predict the response variable mv (median value of homes in thousands of 1970 US dollars) from subsets of four of the explanatory variables, as described in Week 6 Assignment Data: Boston Housing Study. Additionally, the mean-square error, AIC, and BIC were computed as well.
## Implementation
Two key packages used in both implementations are "github.com/gonum/stat/combin" and "github.com/sajari/regression". combin was utilized to create the variable combinations of 4 variables per combination. sajari/regression served as the foundation for the linear regression models.

Each program starts off with some data exploration prior to running the linear regression program. The two programs also share similar functions, including getIndex, generateCombinations, computeMSE, computeAIC, and computeBIC.
## Without vs With Concurrency
One of the main challenges I found when implementing the program with concurrency was to make sure that the outputs were correctly grouped together for each run. My first few implementations would output the all possible combiantions at once, followed by the formulas with MSE, AIC, and BIC outputs dispersed throughout. To work around this, one of the biggest differences in the program utilizing concurrency makes use of struct to store a single iteration's combination, formula, MSE, AIC, and BIC. The output then occurs all at once, instead of printing line by line as information is calculated like the non-concurrent program.
## Run Times & Performance
#### Without concurrency: 
- real: 1:45.64 (~105 seconds)
- user: 1.09 s
- sys: 0.23 s

#### With concurrency: 
- real: 47.624 s
- user: 1.38 s
- sys: 0.47 s

Overall, the program utilizing concurrency performed much better as it cut the performance time by over a minute.
