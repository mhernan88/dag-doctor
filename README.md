# DAG Doctor

**Problem**: Frequently in data engineering (and other applications), elements pass through a Directed Acyclic Graph. In development (and occassionally in production)
things break in a graph. Sometimes it is quite difficult to find where exactly something has broken. In that case, a developer needs to search log files, and
run pipelines in debug mode until they find the source of the problem. This search process can be lengthy, especially if there are long-running nodes in the graph.

The first defense against bugs in a pipeline is a good set of unit and integration tests. Beyond that a data quality layer also helps significantly with reducing the
likelihood of bugs. But, edge cases do pop up that cause unexpected issues. In those cases, DAG Doctor can be of use.

**Solution**: DAG Doctor traverses a Directed Acyclic Graph and recommends optimal graph nodes to test in order to, as quickly as possible, find the source of a bug
(in a similar way that binary search does for an ordered list).
