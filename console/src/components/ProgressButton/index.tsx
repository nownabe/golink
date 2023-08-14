import { Box, Button, ButtonProps, CircularProgress } from "@mui/material";

type Props = {
  loading: boolean;
  disabled?: boolean;
  color?: ButtonProps["color"];
  onClick: () => void;
  children: React.ReactNode;
};

export default function ProgressButton({
  loading,
  disabled,
  color,
  onClick,
  children,
}: Props) {
  return (
    <Box sx={{ display: "flex" }}>
      <Box sx={{ position: "relative" }}>
        <Button
          variant="contained"
          disabled={loading || disabled}
          color={color}
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
