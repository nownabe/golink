"use client";

import { Box, Chip } from "@mui/material";
import { SyntheticEvent, useCallback } from "react";

type Props = {
  owners: string[];
};

export default function DeletableOwners({ owners }: Props) {
  const onDelete = useCallback((email: string) => {
    return () => {
      alert(`delete ${email}`);
    };
  }, []);

  return (
    <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
      {owners.map((owner) => (
        <Chip key={owner} label={owner} onDelete={onDelete(owner)} />
      ))}
    </Box>
  );
}
