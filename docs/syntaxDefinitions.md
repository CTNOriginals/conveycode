# convcode Syntax Definitions

This document contains some rules that will apply to the syntax of convcode. Though, this is an ongoing process and will be updated throughout development as needed.

## Assignment
<!-- - When defining a new variable in any context, the usage of `:=` is required, otherwise, when assigning a value to an already existing variable, the usage of `=` is required instead of `:=` -->
- When defining a new variable in any context, it is required to prefix it with `var`, otherwise, when assigning a value to an already existing variable, dont use a prefix at all