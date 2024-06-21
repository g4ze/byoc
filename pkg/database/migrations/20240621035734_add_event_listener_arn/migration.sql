/*
  Warnings:

  - You are about to drop the column `DesiredCount` on the `Service` table. All the data in the column will be lost.
  - You are about to drop the column `Slug` on the `Service` table. All the data in the column will be lost.
  - A unique constraint covering the columns `[slug]` on the table `Service` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `desiredCount` to the `Service` table without a default value. This is not possible if the table is not empty.
  - Added the required column `eventListenerARN` to the `Service` table without a default value. This is not possible if the table is not empty.
  - Added the required column `slug` to the `Service` table without a default value. This is not possible if the table is not empty.

*/
-- AlterTable
ALTER TABLE "Service" DROP COLUMN "DesiredCount",
DROP COLUMN "Slug",
ADD COLUMN     "desiredCount" INTEGER NOT NULL,
ADD COLUMN     "eventListenerARN" TEXT NOT NULL,
ADD COLUMN     "slug" TEXT NOT NULL;

-- CreateIndex
CREATE UNIQUE INDEX "Service_slug_key" ON "Service"("slug");
