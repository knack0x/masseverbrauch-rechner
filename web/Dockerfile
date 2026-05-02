FROM php:8.2-apache

RUN apt-get update && apt-get install -y \
    libcurl4-openssl-dev \
    && docker-php-ext-install curl

RUN a2enmod rewrite

COPY web/ /var/www/html/
COPY VERSION /var/www/html/VERSION

EXPOSE 80
