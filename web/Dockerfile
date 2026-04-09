FROM node:24-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# ARG REACT_APP_USER_API_URL
# ARG REACT_APP_VIDEO_API_URL

# ENV REACT_APP_USER_API_URL=${REACT_APP_USER_API_URL}
# ENV REACT_APP_VIDEO_API_URL=${REACT_APP_VIDEO_API_URL}

# Build the React app with environment variables
RUN npm run build

FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
