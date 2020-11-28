const useFileImporter = <T>(
  accept: string[],
  cb: (
    filename: string,
    content: string | ArrayBuffer | null | undefined,
    x?: T
  ) => void
) => (x?: T) => {
  const input = document.createElement("input");
  input.type = "file";
  input.accept = accept.join(",");

  input.onchange = (e) => {
    console.log("file changed");

    const file = (e.target as HTMLInputElement)?.files?.[0];
    if (!file) return;

    const reader = new FileReader();
    reader.readAsArrayBuffer(file);

    reader.onload = (readerEvent) => {
      cb(file.name, readerEvent.target?.result, x);
    };

    input.parentNode?.removeChild(input);
  };
  input.oncancel = () => {
    input.parentNode?.removeChild(input);
  };

  input.click();
};

export default useFileImporter;
