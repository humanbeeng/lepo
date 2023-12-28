# Callgraph creation

- I could just have method block as the caller for any call, but if I want to have specific references (go to references) then I would need
  the exact location where the call happened. So I have to manually go inside all kinds of closure where function call is possible.

Bulk imports from CSV or JSON

- Create two files, one for all nodes, and one for all relations.
- First create all nodes and then create relations.

TODO for Jan:

- [ ] Prototype graph creation using json or csv
- [ ] Method calls

### AST

- Anon functions are represented by ast.FuncLit
- function calls within same packages are not represented by selector expr rather by identifier itself

### Unknowns

- [ ] How to generate callgraph for internal functions (stdlib functions)

Generate TS/JS callgraph
[Link](https://github.com/cs-au-dk/jelly)
