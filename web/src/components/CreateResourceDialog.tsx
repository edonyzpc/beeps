import { Button, IconButton, List, ListItem, Typography } from "@mui/joy";
import React, { useRef, useState } from "react";
import { toast } from "react-hot-toast";
import { useResourceStore } from "@/store/v1";
import { Resource } from "@/types/proto/api/v1/resource_service";
import { useTranslate } from "@/utils/i18n";
import { generateDialog } from "./Dialog";
import Icon from "./Icon";

interface Props extends DialogProps {
  onCancel?: () => void;
  onConfirm?: (resourceList: Resource[]) => void;
}

interface State {
  uploadingFlag: boolean;
}

const CreateResourceDialog: React.FC<Props> = (props: Props) => {
  const t = useTranslate();
  const { destroy, onCancel, onConfirm } = props;
  const resourceStore = useResourceStore();
  const [state, setState] = useState<State>({
    uploadingFlag: false,
  });
  const [fileList, setFileList] = useState<File[]>([]);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleReorderFileList = (fileName: string, direction: "up" | "down") => {
    const fileIndex = fileList.findIndex((file) => file.name === fileName);
    if (fileIndex === -1) {
      return;
    }

    const newFileList = [...fileList];

    if (direction === "up") {
      if (fileIndex === 0) {
        return;
      }
      const temp = newFileList[fileIndex - 1];
      newFileList[fileIndex - 1] = newFileList[fileIndex];
      newFileList[fileIndex] = temp;
    } else if (direction === "down") {
      if (fileIndex === fileList.length - 1) {
        return;
      }
      const temp = newFileList[fileIndex + 1];
      newFileList[fileIndex + 1] = newFileList[fileIndex];
      newFileList[fileIndex] = temp;
    }

    setFileList(newFileList);
  };

  const handleCloseDialog = () => {
    if (onCancel) {
      onCancel();
    }
    destroy();
  };

  const handleFileInputChange = async () => {
    if (!fileInputRef.current || !fileInputRef.current.files) {
      return;
    }

    const files: File[] = [];
    for (const file of fileInputRef.current.files) {
      files.push(file);
    }
    setFileList(files);
  };

  const allowConfirmAction = () => {
    if (!fileInputRef.current || !fileInputRef.current.files || fileInputRef.current.files.length === 0) {
      return false;
    }
    return true;
  };

  const handleConfirmBtnClick = async () => {
    if (state.uploadingFlag) {
      return;
    }

    setState((state) => {
      return {
        ...state,
        uploadingFlag: true,
      };
    });

    const createdResourceList: Resource[] = [];
    try {
      if (!fileInputRef.current || !fileInputRef.current.files) {
        return;
      }
      const filesOnInput = Array.from(fileInputRef.current.files);
      for (const file of fileList) {
        const fileOnInput = filesOnInput.find((fileOnInput) => fileOnInput.name === file.name);
        if (!fileOnInput) {
          continue;
        }
        const { name: filename, size, type } = file;
        const buffer = new Uint8Array(await file.arrayBuffer());
        const resource = await resourceStore.createResource({
          resource: Resource.fromPartial({
            filename,
            size,
            type,
            content: buffer,
          }),
        });
        createdResourceList.push(resource);
      }
    } catch (error: any) {
      console.error(error);
      toast.error(error.details);
    }

    if (onConfirm) {
      onConfirm(createdResourceList);
    }
    destroy();
  };

  return (
    <>
      <div className="dialog-header-container">
        <p className="title-text">{t("resource.create-dialog.title")}</p>
        <IconButton size="sm" onClick={handleCloseDialog}>
          <Icon.X className="w-5 h-auto" />
        </IconButton>
      </div>
      <div className="dialog-content-container !w-80">
        <div className="w-full relative bg-blue-50 dark:bg-zinc-900 rounded-md border-dashed border-2 dark:border-zinc-700 flex flex-row justify-center items-center py-8 hover:opacity-90">
          <label htmlFor="files" className="p-2 px-4 text-sm text-white cursor-pointer bg-blue-500 block rounded hover:opacity-80">
            {t("resource.create-dialog.local-file.choose")}
          </label>
          <input
            className="absolute inset-0 w-full h-full opacity-0"
            ref={fileInputRef}
            onChange={handleFileInputChange}
            type="file"
            id="files"
            multiple={true}
            accept="*"
          />
        </div>
        <List size="sm" sx={{ width: "100%" }}>
          {fileList.map((file, index) => (
            <ListItem key={file.name} className="flex justify-between">
              <Typography noWrap>{file.name}</Typography>
              <div className="flex gap-1">
                <button
                  onClick={() => {
                    handleReorderFileList(file.name, "up");
                  }}
                  disabled={index === 0}
                  className="disabled:opacity-50"
                >
                  <Icon.ArrowUp className="w-4 h-4" />
                </button>
                <button
                  onClick={() => {
                    handleReorderFileList(file.name, "down");
                  }}
                  disabled={index === fileList.length - 1}
                  className="disabled:opacity-50"
                >
                  <Icon.ArrowDown className="w-4 h-4" />
                </button>
              </div>
            </ListItem>
          ))}
        </List>

        <div className="mt-2 w-full flex flex-row justify-end items-center space-x-1">
          <Button variant="plain" color="neutral" onClick={handleCloseDialog}>
            {t("common.cancel")}
          </Button>
          <Button onClick={handleConfirmBtnClick} loading={state.uploadingFlag} disabled={!allowConfirmAction()}>
            {t("common.create")}
          </Button>
        </div>
      </div>
    </>
  );
};

function showCreateResourceDialog(props: Omit<Props, "destroy">) {
  generateDialog<Props>(
    {
      dialogName: "create-resource-dialog",
    },
    CreateResourceDialog,
    props,
  );
}

export default showCreateResourceDialog;
