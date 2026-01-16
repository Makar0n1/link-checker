"use client";

import { useState, useEffect, useCallback } from "react";
import { Input } from "@/components/ui/input";

interface EditableCellProps {
  value: string;
  onSave: (value: string) => void;
}

export function EditableCell({ value: initialValue, onSave }: EditableCellProps) {
  const [value, setValue] = useState(initialValue);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    setValue(initialValue);
  }, [initialValue]);

  const handleBlur = useCallback(() => {
    setIsEditing(false);
    if (value !== initialValue) {
      onSave(value);
    }
  }, [value, initialValue, onSave]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent) => {
      if (e.key === "Enter") {
        handleBlur();
      }
      if (e.key === "Escape") {
        setValue(initialValue);
        setIsEditing(false);
      }
    },
    [handleBlur, initialValue]
  );

  if (isEditing) {
    return (
      <Input
        value={value}
        onChange={(e) => setValue(e.target.value)}
        onBlur={handleBlur}
        onKeyDown={handleKeyDown}
        autoFocus
        className="h-8 w-full min-w-[100px]"
      />
    );
  }

  return (
    <div
      className="cursor-pointer truncate rounded px-2 py-1 hover:bg-muted"
      onClick={() => setIsEditing(true)}
      title="Click to edit"
    >
      {value || <span className="text-muted-foreground">-</span>}
    </div>
  );
}
