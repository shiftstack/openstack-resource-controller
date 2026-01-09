# Update Trunk

## Step 00

Create a Trunk using only mandatory fields.

## Step 01

Update all mutable fields, except for `AdminStateUp`, since neutron disallow operations on disabled ports.

## Step 02

Update `AdminStateUp`.

## Step 03

Revert the resource to its original value and verify that the resulting object matches its state when first created.

## Reference

https://k-orc.cloud/development/writing-tests/#update
