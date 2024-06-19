/*
  Warnings:

  - Added the required column `cluster` to the `Service` table without a default value. This is not possible if the table is not empty.
  - Added the required column `image` to the `Service` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Service" ADD COLUMN     "cluster" TEXT NOT NULL,
ADD COLUMN     "image" TEXT NOT NULL;
