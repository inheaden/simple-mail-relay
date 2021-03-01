FROM alpine:3.12.1

ENV USER=appuser
ENV UID=2001
ENV GID=2001

RUN addgroup \
  --gid "${GID}" \
  "${USER}"

RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/home/${USER}" \
  --ingroup "${USER}" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"

WORKDIR /home/${USER}

# Add a volume pointing to /tmp
VOLUME /tmp

# Copy our static executable.
COPY ./service service

# empty .env file
RUN touch .env
RUN chown -R $UID:$GID /home/${USER} & chown -R $UID:$GID  /tmp
RUN chmod +x ./service

USER ${USER}

ENTRYPOINT ["/home/appuser/service"]