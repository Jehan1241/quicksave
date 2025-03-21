import { toast } from "@/hooks/use-toast"; // Ensure this is correctly imported

export const showErrorToast = (title: string, description: string) => {
  toast({
    variant: "destructive",
    title: title,
    description: description,
  });
};
