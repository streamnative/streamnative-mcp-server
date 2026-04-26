# Lane: Fix

This lane fixes a defect found in the Pulsar resource implementation.

Focus on:

- reproducing the defect with a focused test
- fixing the smallest affected surface
- preserving existing URI compatibility unless the URI is unsafe or invalid
- updating docs only when user-visible behavior changes

If the defect reveals a broader design issue, fix the immediate bug and record
the larger follow-up in `next_suggestion`.
