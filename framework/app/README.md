Application Container
=====================

# 1. Feature

* [X] Register/Retrieve object to/from container
* [X] Simple dependency injection
* [ ] Simple object lifecycle management

Note:

* Any unsatisfied operation provided by the container will cause panic.
    * This guarantees the completeness of the objects in the container for the whole application
* Have to meet all requirements in order to use the features provided by the container
    * Refer to the comments on each method for the details
* Should not use circle dependency in dependency injection
    * Otherwise, it will cause unexpected behavior
* Initialization/Destruction order will be automatically calculated inside the container according to dependency
  injection graph.
