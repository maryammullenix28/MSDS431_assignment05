# Exploring Concurrency
This assignment asked students to explore using Go's concurrency features (channels and goroutines in particular) on a program that fits all possible linear regretions for the Boston Housing Study. The code considers only main-effect models, with no-interactions and no transformations of the explanatory variables. The models predict the response variable mv (median value of homes in thousands of 1970 US dollars) from subsets of four of the explanatory variables, as described in Week 6 Assignment Data: Boston Housing Study. Additionally, the mean-square error, AIC, and BIC were computed as well.
## Implementation
Two key packages used in both implementations are "github.com/gonum/stat/combin" and "github.com/sajari/regression". combin was utilized to create the variable combinations of 4 variables per combination. sajari/regression served as the foundation for the linear regression models.
## Without Concurrency
