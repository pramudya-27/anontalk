FROM redis:alpine

# Optional: Copy custom configuration if needed
# COPY redis.conf /usr/local/etc/redis/redis.conf
# CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]

CMD [ "redis-server" ]
EXPOSE 6379
