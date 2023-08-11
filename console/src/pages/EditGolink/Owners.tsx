import { Box, Chip, TextField, Typography } from "@mui/material";
import Grid from "@mui/material/Unstable_Grid2";

import ProgressButton from "@/components/ProgressButton";
import Snackbar from "@/components/Snackbar";

import { useOwners } from "./hooks";
import { Golink } from "@/gen/golink/v1/golink_pb";

type Props = {
  golink: Golink;
};

export default function Owners({ golink }: Props) {
  const {
    owners,
    removeOwner,
    openRemoveSuccess,
    onRemoveSuccessClose,
    removeError,
    onRemoveErrorClose,
    addOwner,
    addRef,
    adding,
    openAddSuccess,
    onAddSuccessClose,
    addError,
    onAddErrorClose,
  } = useOwners(golink);

  return (
    <>
      <Grid xs={12}>
        <Typography variant="h6" component="h3">
          Owners
        </Typography>
      </Grid>
      <Grid xs={12}>
        <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
          {owners.map((owner) => (
            <Chip
              key={owner}
              label={owner}
              onDelete={owners.length !== 1 ? removeOwner(owner) : undefined}
            />
          ))}
          <Snackbar
            open={!!removeError}
            severity="error"
            onClose={onRemoveErrorClose}
          >
            {removeError || ""}
          </Snackbar>
          <Snackbar
            open={openRemoveSuccess}
            severity="success"
            onClose={onRemoveSuccessClose}
          >
            Successfully deleted.
          </Snackbar>
        </Box>
      </Grid>
      <Grid xs={12}>
        <Typography variant="h6" component="h3">
          Add Owner
        </Typography>
      </Grid>
      <Grid xs={12}>
        <TextField
          label="email"
          inputRef={addRef}
          fullWidth
          placeholder="new-owner@example.com"
        />
      </Grid>
      <Grid xs={2}>
        <ProgressButton loading={adding} onClick={addOwner}>
          Add Owner
        </ProgressButton>
      </Grid>
      <Snackbar open={!!addError} severity="error" onClose={onAddErrorClose}>
        {addError || ""}
      </Snackbar>
      <Snackbar
        open={openAddSuccess}
        severity="success"
        onClose={onAddSuccessClose}
      >
        Successfully added.
      </Snackbar>
    </>
  );
}
