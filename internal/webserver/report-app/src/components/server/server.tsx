import { useToast } from "use-toast-mui";

// define host from window.location.hostname, or if it is localhost, use "http://hayden.lavina.tech"

let hosVar = window.location.hostname;
if (hosVar === "localhost") {
  hosVar = "mason.lavina.uz";
}

export const host = "https://" + hosVar;

export const fetchServer = (
  input: RequestInfo | URL,
  init?: RequestInit
): Promise<Response> => {
  init = {
    ...init,
    headers: {
      ...init?.headers,
      "Content-Type": "application/json",
    },
  };
  return fetch(host + input, init);
};

type SrvRes = {
  data: any;
  message: string;
  isOk: boolean;
};

export const useFetchServer = (): (<DataType>(
  input: RequestInfo | URL,
  init?: RequestInit
) => Promise<DataType>) => {
  const toast = useToast();
  return <DataType,>(
    input: RequestInfo | URL,
    init?: RequestInit
  ): Promise<DataType> => {
    init = {
      ...init,
      headers: {
        ...init?.headers,
        "Content-Type": "application/json",
      },
    };
    let msg = "";
    return new Promise<DataType>((resolve, reject) => {
      fetch(host + input, init)
        .then((res: Response) => {
          res.text().then((body) => {
            if (body === "" && res.ok) {
              resolve(undefined as DataType);
              return;
            }

            if (body === "" && !res.ok) {
              msg = "Unable to fetch " + input + ": ";
              toast.error(msg);
              reject(msg);
              return;
            }

            if (body !== "" && res.ok) {
              try {
                const data = JSON.parse(body);
                resolve(data.data as DataType);
              } catch (error) {
                toast.error("Unable to read result of " + input);
                reject(error);
              }
            }

            if (body !== "" && !res.ok) {
              try {
                const data = JSON.parse(body);
                msg = "Unable to fetch " + input + ": " + data.message;
                toast.error(msg);
                reject(msg);
              } catch (error) {
                toast.error("Unable to read result of " + input);
                msg = "Unable to fetch " + input + ": ";
                toast.error(msg);
                reject(msg);
              }
            }
          });
        })
        .catch((res) => {
          toast.error("Request failed");
          reject(res);
        });
    });
  };
};
