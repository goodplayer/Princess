Application Container
=====================

# 1. Feature

* [X] Register/Retrieve object to/from container
* [X] Simple dependency injection
* [ ] Simple object lifecycle management
* [ ] Advanced dependency injection by specifying provider name
* [ ] Advanced dependency injection by interface
* [X] Support singleton
* [ ] Support prototype - may not support in short term
* [X] Support ptr to struct for DI
* [ ] Support struct for DI - may not support in short term
* [ ] Custom DI: custom injection method
* [ ] Custom DI: custom turn on/off lifecycle management

Note:

* Any unsatisfied operation provided by the container will cause panic.
    * This guarantees the completeness of the objects in the container for the whole application
* Have to meet all requirements in order to use the features provided by the container
    * Refer to the comments on each method for the details
* Should not use circle dependency in dependency injection
    * Otherwise, it will cause unexpected behavior
* Initialization/Destruction order will be automatically calculated inside the container according to dependency
  injection graph.

# 2. Discussion

* For dependency injection by name
    * The solution should be directly supporting name for injection
    * An alternative is to use a sub container which contains the new providers and same requester. In this case, sub
      container will construct partial of the dependency while main container construct the rest and controls the
      lifecycle of the objects.
    * Looks like the second approach is weird but works fine as well.
* Using ptr to a struct or struct directly as object injected in the container
    * Using struct/factory that produces struct, it is natural and suggested to be the best practice in dependency
      injection. In common cases, user would prefer to use a NewObject function which consumes dependencies and produces
      desired object.
    * Using ptr to a struct for dependency injection. It will limit the use cases for a user to provide a ptr only. But
      it will provide better reliability when handling value copy issue. And it will also allow user to use an object
      right after dependency inject without having to get the object from the container(which is not suggested in
      practice).

