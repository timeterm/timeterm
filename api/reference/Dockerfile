FROM nginx:alpine

COPY /index.html /usr/share/nginx/html
COPY /openapi.yaml /usr/share/nginx/html
COPY nginx/nginx.conf /etc/nginx/

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

