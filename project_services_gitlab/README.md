## Installing the Project Service (GitLab)

The project service can be installed by coping the file `learning_bot_service.rb` from
this folder to the GitLab-CE source at `app/models/project_services/[learning_bot_service.rb]`.  

After copying the file, you will have to change the `learning_bot_service_url` constant
to be the bot's instance URL (including port, if non-default).  

Then make sure to include the service models in `app/models/project.rb` and
`spec/models/project_spec.rb` to make sure that GitLab recognises and loads
the new service.
