MRuby::Build.new do |conf|
  # load specific toolchain settings

  # Gets set by the VS command prompts.
  if ENV['VisualStudioVersion'] || ENV['VSINSTALLDIR']
    toolchain :visualcpp
  else
    toolchain :gcc
  end

  if ENV['CI']
    # Turn on `enable_debug` for better debugging
    enable_debug
  end

  gem core: 'mruby-error'
end
