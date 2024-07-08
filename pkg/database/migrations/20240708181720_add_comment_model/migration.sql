-- CreateTable
CREATE TABLE "User" (
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "userName" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "password" TEXT NOT NULL,
    "images" TEXT,
    "cluster" TEXT NOT NULL,

    CONSTRAINT "User_pkey" PRIMARY KEY ("userName")
);

-- CreateTable
CREATE TABLE "Service" (
    "id" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "name" TEXT NOT NULL,
    "arn" TEXT NOT NULL,
    "taskFamily" TEXT NOT NULL,
    "loadBalancerARN" TEXT NOT NULL,
    "targetGroupARN" TEXT NOT NULL,
    "loadbalancerDNS" TEXT NOT NULL,
    "desiredCount" INTEGER NOT NULL,
    "cluster" TEXT NOT NULL,
    "image" TEXT NOT NULL,
    "slug" TEXT NOT NULL,
    "eventListenerARN" TEXT NOT NULL,
    "userName" TEXT NOT NULL,
    "logs" TEXT,
    "deploymentName" TEXT NOT NULL DEFAULT 'default',

    CONSTRAINT "Service_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "User_userName_key" ON "User"("userName");

-- CreateIndex
CREATE UNIQUE INDEX "User_email_key" ON "User"("email");

-- CreateIndex
CREATE UNIQUE INDEX "Service_slug_key" ON "Service"("slug");

-- AddForeignKey
ALTER TABLE "Service" ADD CONSTRAINT "Service_userName_fkey" FOREIGN KEY ("userName") REFERENCES "User"("userName") ON DELETE RESTRICT ON UPDATE CASCADE;
