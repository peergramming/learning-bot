# Based on MockCiService
class LearningBotService < CiService # this service follows on the same aspects of CI
  ALLOWED_STATES = %w[failed canceled running pending success success-with-warnings skipped not_found].freeze

  def title
    'Learning Bot'
  end

  def description
    'Code checking bot'
  end

  def self.to_param
    'learning_bot'
  end

  def learning_bot_service_url # Constant to reduce configuration for each student
    'http://gitlab-student'
  end

  # Return complete url to build page
  #
  def build_page(sha, ref)
    Gitlab::Utils.append_path(
      learning_bot_service_url,
      "#{project.namespace.path}/#{project.path}/report/#{sha}")
  end

  # Return string with build status or :error symbol
  #
  # Allowed states: 'success', 'failed', 'running', 'pending', 'skipped'
  #
  def commit_status(sha, ref)
    response = Gitlab::HTTP.get(commit_status_path(sha), verify: false)
    read_commit_status(response)
  rescue Errno::ECONNREFUSED
    :error
  end

  def commit_status_path(sha)
    Gitlab::Utils.append_path(
      mock_service_url,
      "#{project.namespace.path}/#{project.path}/status/#{sha}.json")
  end

  def read_commit_status(response)
    return :error unless response.code == 200 || response.code == 404

    status = if response.code == 404
               'pending'
             else
               response['status']
             end

    if status.present? && ALLOWED_STATES.include?(status)
      status
    else
      :error
    end
  end

  def can_test?
    false
  end
end
