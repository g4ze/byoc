/*
  Warnings:

  - You are about to drop the column `proxyURL` on the `Service` table. All the data in the column will be lost.
  - Added the required column `Slug` to the `Service` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Service" DROP COLUMN "proxyURL",
ADD COLUMN     "Slug" TEXT NOT NULL;
