# DAG Doctor

Quickly pinpoint the source of silent data pipeline faults.

**Scenario**: You and your team have created a large data pipeline. After months of releases and enhancements, you get a call from a consumer of your data pipeline stating that the data looks wrong - and they don't know when this issue first popped up. You have unit tests in place and a data cleansing layer in place, so you're not sure where to begin to find the fault.

**Old Solution**: You run through the (potentially hundreds of) nodes in your DAG (directed acyclic graph) and find the fault after two weeks of searching. Turns out it was due to an edge case that was not caught by unit testing.

**New Solution**: You start up DAG Doctor. DAG Doctor helps you intelligently bisect your graph into smaller and smaller search spaces until your offending node is located. Think of it as binary search, but for a DAG.
