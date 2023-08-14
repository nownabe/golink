import { Box, Button, CircularProgress } from "@mui/material";

type Props = {
  loading: boolean;
  disabled?: boolean;
  onClick: () => void;
  children: React.ReactNode;
};

export default function ProgressButton({
  loading,
  disabled,
  onClick,
  children,
}: Props) {
  return (
    <Box sx={{ display: "flex" }}>
      <Box sx={{ position: "relative" }}>
        <Button
          variant="contained"
          disabled={loading || disabled}
          onClick={onClick}
        >
          {children}
        </Button>
        {loading && (
          <CircularProgress
            size={24}
            sx={{
              position: "absolute",
              top: "50%",
              left: "50%",
              marginTop: "-12px",
              marginLeft: "-12px",
            }}
          />
        )}
      </Box>
    </Box>
  );
}
