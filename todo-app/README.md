# Todo app backend

First navigate to the `todo-app` directory:

```shell
cd todo-app
```

To deploy in staging environment, run:

```shell
kubectl apply -f ./kubernetes/overlays/staging/
```

or in production environment, run:

```shell
kubectl apply -f ./kubernetes/overlays/production/
```

Note: The application is highly dependent on Google Cloud Artifact Registry, so deploying from here might not work.

# Exercise 3.9 DBaaS vs DIY

### DBaaS (Database as a Service)

**Pros:**

- **Easy Setup:** You can set up a database in minutes, often with just a few clicks.
- **Predictable cost:** Providers usually offer multiple pricing plans.
- **Maintenance:** The provider takes care of updates, scaling, and security patches for you.
- **Backups:** Backups are usually automatic and restoring data is straightforward.
- **Reliability:** High availability and disaster recovery are built in.

**Cons:**

- **Potential Expensiveness:**
- **Limitations:** You can’t tweak every setting or access the underlying infrastructure.
- **Vendor lock-in:** Migrating data to another provider might be challenging.

### DIY (Self-Managed Database)

**Pros:**

- **Full control:** You decide how everything is set up and configured.
- **Potential savings:** If you already have hardware, running your own database can be cheaper for big projects.
- **Flexibility:** You can use any backup, scaling, or monitoring tools you like.

**Cons:**

- **Setup:** You’ll need to set up hardware, install software, and handle networking and security yourself. If you
  don't have someone who knows how to do this, it can be a lot of work.
- **Hardware Cost:** Servers are usually quite expensive.
- **Ongoing maintenance:** Updates, scaling, monitoring, and fixing bugs and other issues fall on you.
- **Manual backups:** You have to set up and regularly test your own backup and restore processes.
- **Reliability takes effort:** High availability and disaster recovery require extra planning and resources.

### How easy are backups?

- **DBaaS:** Backups are usually handled for you. You can restore data easily through a dashboard or API.
- **DIY:** You need to set up your own backup scripts or use third-party tools. Restoring data may take more time and
  testing.

---
