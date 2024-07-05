FROM nginx
LABEL authors="noroom113"
RUN rm /etc/nginx/conf.d/default.conf
COPY content /usr/share/nginx/html
COPY conf /etc/nginx/conf.d
VOLUME /usr/share/nginx/html
VOLUME /etc/nginx
EXPOSE 80:80