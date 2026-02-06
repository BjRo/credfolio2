import { FileText, Plus, Upload } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ProfileActionsProps {
  onAddReference?: () => void;
  onExport?: () => void;
  onUploadAnother?: () => void;
}

interface ActionItem {
  label: string;
  icon: React.ComponentType<{ className?: string }>;
  onClick: () => void;
  variant: "outline" | "ghost";
}

function buildActions({
  onAddReference,
  onExport,
  onUploadAnother,
}: ProfileActionsProps): ActionItem[] {
  const actions: ActionItem[] = [];

  if (onAddReference) {
    actions.push({
      label: "Add Reference Letter",
      icon: Plus,
      onClick: onAddReference,
      variant: "outline",
    });
  }
  if (onExport) {
    actions.push({ label: "Export PDF", icon: FileText, onClick: onExport, variant: "outline" });
  }
  if (onUploadAnother) {
    actions.push({
      label: "Upload Another Resume",
      icon: Upload,
      onClick: onUploadAnother,
      variant: "ghost",
    });
  }

  return actions;
}

export function ProfileActions(props: ProfileActionsProps) {
  const actions = buildActions(props);

  if (actions.length === 0) return null;

  return (
    <div className="lg:hidden bg-card border rounded-lg p-6">
      <div className="flex flex-wrap justify-center gap-4">
        {actions.map((action) => (
          <Button key={action.label} variant={action.variant} onClick={action.onClick}>
            <action.icon className="w-4 h-4 mr-2" aria-hidden="true" />
            {action.label}
          </Button>
        ))}
      </div>
    </div>
  );
}

export function ProfileActionsBar(props: ProfileActionsProps) {
  const actions = buildActions(props);

  if (actions.length === 0) return null;

  return (
    <div className="hidden lg:flex flex-col gap-2 bg-card border rounded-lg p-2">
      {actions.map((action) => (
        <div key={action.label} className="group relative">
          <Button
            variant={action.variant}
            size="icon"
            onClick={action.onClick}
            aria-label={action.label}
            className="hover:scale-110"
          >
            <action.icon className="w-4 h-4" aria-hidden="true" />
          </Button>
          <span
            role="tooltip"
            className="pointer-events-none absolute right-full top-1/2 -translate-y-1/2 mr-2 whitespace-nowrap rounded-md bg-popover text-popover-foreground border px-2 py-1 text-xs font-medium shadow-md opacity-0 scale-95 group-hover:opacity-100 group-hover:scale-100 group-focus-within:opacity-100 group-focus-within:scale-100 transition-all duration-150"
          >
            {action.label}
          </span>
        </div>
      ))}
    </div>
  );
}
