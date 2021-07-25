import React from "react";
import {Store} from "pullstate";
import {Alert, ToastContainer} from "react-bootstrap";
import {Variant} from "react-bootstrap/types";

import "./Toasts.scss";

export type ToastType = {
  title?: string;
  message?: string;
  type?: Variant;
}

export const ToastStore = new Store<{ toasts: ToastType[] }>({
  toasts: []
});

export function addToast(toast: ToastType): void {
  ToastStore.update((s) => {
    s.toasts.push(toast);
  });
  setTimeout(removeToast.bind(0), 15000);
}

function removeToast(idx: number) {
  ToastStore.update((s) => {
    s.toasts.splice(idx, 1);
  });
}

export function ToastComponent(): JSX.Element {
  const toasts = ToastStore.useState(s => s.toasts);
  return (
    <ToastContainer className={"toast-container p-3"} position={"top-end"}>
      {toasts.map((toast, idx) =>
        <Alert key={idx} onClose={removeToast.bind(idx)} variant={toast.type || "info"} dismissible>
          {toast.title ? <Alert.Heading>{toast.title}</Alert.Heading> : []}
          {toast.message ? <span className={"me-auto"}>{toast.message}</span> : []}
        </Alert>
      )}
    </ToastContainer>
  );
}