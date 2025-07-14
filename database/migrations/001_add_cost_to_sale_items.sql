-- Migration: Add cost field to sale_items table
-- Date: 2025-07-14

-- Add cost column to sale_items table
ALTER TABLE sale_items ADD COLUMN IF NOT EXISTS cost DECIMAL(10,2) NOT NULL DEFAULT 0;

-- Update existing sale_items with product cost at time of creation
UPDATE sale_items 
SET cost = products.cost 
FROM products 
WHERE sale_items.product_id = products.id 
AND sale_items.cost = 0;