# Closed loop demos
In this repository, demos presenting the operation of control loop framework based on Kubernetes operator pattern is presented. The basic assumption behind the framework is that a control loop follows the architecture of autonomous systems where several components can operate in concert to provide the desired functionality of the loop. Each component of the loop can be implemented as a separate Kubernetes operator (i.e., custom resource and respective custom controller) and the user can define loop architecture (the set of components involved in loop instance) and to some degree detail the logic of the loop in a declarative way using Kubernetes custom resources. More on the philosophy and the assumptions behind the framework can be found in the reports that are referenced in user guides (Demo1.md, Demo2.md documents) available in this repository. Those guides describe, respectively:

**[Demo1.md](./docs/Demo1.md?ref_type=heads)**: Demo-1 - guide for the demo summarizing the work done in 2023. Two cases are described in this document. The first one presents a single reactive control loop operating in isolation (simple case). The second one presents hierarchical setup of two interworking control loops (reactive control loop and deliberative control loop where deliberative control loop monitors and adjusts the operating parameters of the reactive control loop). A distinctive feature of '2023 framework used in Demo-1 is the decision making logic of loop components being entirely hardcoded in respective custom controllers. In such a case, any change of the logic requires code changes and building new image of respective custom controller. This limits the degree of declarativeness in defining the desired operation of the loop by the user. The document also contains a detailed guide explaining how to install and run the whole environment and deploy control loops in minikube. It is mandatory if one attempts to recreate any of the demos presented in this repo on her/his own.

**[Demo2](./README2.md)** Demo-2 - guide for the demo summarizing the work done in 2024. In this version, the framework from Demo-1 has been enhanced by allowing to delegate the decision making logic of the loop to external applications (through sending queries and receiving responses/decisions). This significantly broadens the possibilities offered to the user to declaratively define the decision taking logic of loop components without the need to recompile the code and build new images of operator containers. In our demo, we use OPA/Rego policy engine as an example of external application playing the role of policy decision point. Mastering the installation of the environment and loop deployment process as outlined in Demo1.md is required to sucessfully recreate Demo-2.
