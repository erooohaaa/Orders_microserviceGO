# 🛒 Order & Payment Microservices

> Educational project in Go --- two microservices using Clean
> Architecture, REST API, and PostgreSQL.

------------------------------------------------------------------------

## 📌 What is this?

I built two small services that work together:

-   **Order Service** --- receives customer orders\
-   **Payment Service** --- processes payments

When a new order is created, the Order Service calls the Payment Service
and asks: "Can we charge this amount?". If yes --- the order becomes
**Paid**. If not --- it is marked as **Failed**.

Each service is a **separate application** with its **own database**.

------------------------------------------------------------------------

## 🔍 How it works

User → Order Service → Payment Service → Order updated → Response to
user

------------------------------------------------------------------------

## 🏗️ Architecture

Clean Architecture means: - Business logic is independent - Database is
separate - HTTP layer is separate

Layers: - domain --- models - usecase --- logic - repository ---
database - transport/http --- handlers

Each service has its own DB.

------------------------------------------------------------------------

## 🚀 How to Run

``` bash
git clone https://github.com/your-username/your-repo.git
cd your-repo
```

Create DB:

``` sql
CREATE DATABASE order_db;
CREATE DATABASE payment_db;
```

Run services:

``` bash
cd Payments && go run ./cmd/payment-service/main.go
cd Orders && go run ./cmd/order-service/main.go
```

------------------------------------------------------------------------

## 🎓 What I learned

-   Clean Architecture basics
-   Microservices communication via REST
-   Working with PostgreSQL in Go
-   Handling failures and timeouts
