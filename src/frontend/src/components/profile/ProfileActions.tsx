import { FileText, Plus, Upload } from "lucide-react";
import { Button } from "@/components/ui/button";

interface ProfileActionsProps {
  onAddReference?: () => void;
  onExport?: () => void;
  onUploadAnother?: () => void;
}

export function ProfileActions({ onAddReference, onExport, onUploadAnother }: ProfileActionsProps) {
  return (
    <div className="bg-card shadow rounded-lg p-6">
      <div className="flex flex-wrap justify-center gap-4">
        {onAddReference && (
          <Button variant="outline" onClick={onAddReference}>
            <Plus className="w-4 h-4 mr-2" aria-hidden="true" />
            Add Reference Letter
          </Button>
        )}
        {onExport && (
          <Button variant="outline" onClick={onExport}>
            <FileText className="w-4 h-4 mr-2" aria-hidden="true" />
            Export PDF
          </Button>
        )}
        {onUploadAnother && (
          <Button variant="ghost" onClick={onUploadAnother}>
            <Upload className="w-4 h-4 mr-2" aria-hidden="true" />
            Upload Another Resume
          </Button>
        )}
      </div>
    </div>
  );
}
