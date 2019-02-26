# Pligos

Pligos allows you to navigate multiple services by making kubernetes
infrastructure configuration scalable. Without pligos, the usual
approach is to create a seperate helm chart for each service. While
this definitely can scale for a small amount of services, maintaining
`deployment, service, ingress, ...` templates for more than 5 helm
charts can be burdensome and error prone.

We observed that services, in it's core, often don't differ that
much. You can find a set of configurations that need to be
individually defined for each service, such as `images, routes,
mounts, ...`, so why not standardize around these configuration types,
while beeing disconnected from the underlying templates? This is why
pligos let's you define these configuration types (`image, route,
...`) and adds a schema language that allows you to compile those
configs into any form necessary.

So, what you will end up with is a small set of helm starters (in
pligos lingua franca they are called flavors) and a pligos
configuration for each service that map to these flavors.

## Schema Compiler

Pligos heavily relies on yaml configurations. In order to compile one
set of configurations into different contexts (for example
CI,dev,prod) pligos comes with it's own simple schema language to
describe the context. The idea is to create a single helm starter that
supports a big set of your services. Services than differ only in
their templating input, or the {values.yaml,dependencies.yaml}
files. Have a look at the examples in the `examples/` directory for a
more comprehensive use-case.

Pligos supports the following basic types: `string, numeric, bool,
object`. Example:

```yaml
# schema.yaml
context:
  pullPolicy: string
  useTLS: bool
  applicationConfig: object
  podInstances: numeric
```

```yaml
# pligos/example/values.yaml -- this is your actual service configuration, define this for all your services individually
contexts:
  dev:
    pullPolicy: Always
    useTLS: true
    applicationConfig:
      fixture:
        user: "John Doe"
        balance: "10$"
    podInstances: 1
```

```yaml
# helmcharts/example/values.yaml -- this file will be generated by pligos
pullPolicy: "Always"
useTLS: true
applicationConfig:
  fixture:
    user: "John Doe"
    balance: "10$"
podInstances: 1
```

As mentioned above pligos allows defining different contexts (dev,
prod, ...). Each context needs to be set under the `contexts` property
in the service configuration. The contexts can be applied to different
schema definitions.

```yaml
# dev/schema.yaml -- the schema definition for the dev context
context:
  environment: string
  srcPath: string
```

```yaml
# prod/schema.yaml -- the schema definition for the prod context
context:
  environment: string
```

```yaml
# pligos/example/values.yaml
contexts:
  dev:
    environment: dev
    srcPath: /home/johndoe/src/
  prod:
    environment: prod
```

```yaml
# helmcharts/dev/values.yaml
environment: dev
srcPath: /home/johndoe/src/
```

```yaml
# helmcharts/prod/values.yaml
environment: prod
```

However the true power of pligos comes through custom types, which can
be instantiated and composed in the service configuration.

```yaml
# schema.yaml
route:
  port: string

container:
  route: route

context:
  container: container
```

```yaml
# pligos/example/values.yaml
route:
 - name: http
   port: 80

container:
 - name: gowebservice
   route: http

contexts:
  dev:
    container: gowebservice
```

```yaml
# helmcharts/example/values.yaml
container:
  route:
    port: 80
```

Notice that in the service configuration the instances are referenced
by their name. `name` is a special property in pligos which is used
for referencering configuration instances (such as the gowebservice
container) and `mapped types`. More to `mapped types` later.

Additionally the language supports the meta types `repeated, mapped,
embedded and embedded mapped` which can be applied to any custom, or
basic types. Let's start with an example for `repeated` instances.

```yaml
# schema.yaml
container:
  name: string
  command: repeated string

context:
  container: repeated container
```

```yaml
# pligos/example/values.yaml
container:
 - name: nginx
   command: ["nginx"]
 - name: php
   command: ["php-fpm"]

contexts:
  dev:
    container: ["nginx", "php"]
```

```yaml
# helmcharts/examples/values.yaml
container:
 - name: nginx
   command:
    - nginx
 - name: php
   command:
    - php-fpm
```

Repeated allows specifying a list of any type. In the example we have
a list of `container`, as well as a property `command` which is a
list of `string`. Next, beside lists, pligos also allows you to define
maps.

```yaml
# schema.yaml
route:
  name: string
  port: string
  containerPort: string

context:
  routes: mapped route
```

```yaml
# pligos/example/values.yaml
route:
 - name: http
   port: 80
   containerPort: 8080

contexts:
  dev:
    routes: ["http"]
```

```yaml
# helmcharts/example/values.yaml
routes:
  http:
    port: 80
    containerPort: 8080
```

Notice that, although the configuration defines an array of routes,
pligos yields a map, as shown in the the output. Maps can be created
using the `mapped` meta type.

Up until this point all custom types appeared under some key in the
output. However, this is not always the desired behavior. In order to
embed the types' properties into the parent `embedded` types can be
used.

```yaml
# schema.yaml
rawValues:
  values: embedded object

context:
  rawValues: embedded rawValues
```

```yaml
# pligos/example/values.yaml
rawValues:
 - name: devenvironment
   values:
     mysql:
       user: testuser
       password: asdf

contexts:
  dev:
    rawValues: devenvironment
```

```yaml
# helmcharts/example/values.yaml
mysql:
  user: testuser
  password: asdf
```

Similarly the instance can be embedded using any arbitrary key using
`embedded mapped` types.

```yaml
# schema.yaml
dependency:
  port: string
  hostname: string

context:
  dependencies: embedded mapped dependency
```

```yaml
# pligos/example/values.yaml
dependency:
 - name: mysql
   port: 3306
   hostname: mysql

contexts:
  dev:
    dependencies: ["mysql"]
```

```yaml
# helmcharts/example/values.yaml
mysql:
  port: 3306
  hostname: mysql
```

Just as with `mapped` types, the name property is used to embedd the
type instance into the parent.
